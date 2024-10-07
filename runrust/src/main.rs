use ini::Ini;
use std::collections::HashMap;
use std::error::Error;
use std::sync::Arc;
use serde::{Serialize, Deserialize};
use actix_web::{web, App, HttpServer, HttpResponse, HttpRequest, Responder, Error as ActixError};
use actix_web::middleware::Logger;
use env_logger::Env;
use chrono::Local;
use std::io::Write;
use once_cell::sync::Lazy;
use std::process::{Command as StdCommand, Stdio};
use std::io::{BufReader, BufRead};
use std::time::{Duration, Instant};
use tokio::time::timeout;
use log::{info, error};
use thiserror::Error;
use clap::Parser;
use rust_embed::RustEmbed;
use mime_guess::from_path;

// Embed frontend files
#[derive(RustEmbed)]
#[folder = "frontend/dist"]
struct FrontendAsset;

type ConfigMap = HashMap<String, String>;

#[derive(Parser, Debug)]
#[clap(author, version, about, long_about = None)]
struct Args {
    #[clap(short, long, default_value = "config.ini")]
    config: String,
}

#[derive(Clone)]
struct Config {
    webconf: ConfigMap,
    paths: ConfigMap,
    commands: HashMap<String, CustomCommand>,
}

#[derive(Clone)]
struct CustomCommand {
    name: Arc<str>,
    value: Arc<str>,
    command: Arc<str>,
    description: Arc<str>,
}

#[derive(Serialize, Deserialize)]
struct ApiCommand {
    name: String,
    value: String,
    command: String,
    description: String,
}

#[derive(Serialize)]
struct Response {
    success: bool,
    message: String,
    details: Option<String>,
}

#[derive(Error, Debug)]
pub enum CommandError {
    #[error("Command execution timed out")]
    Timeout,
    #[error("Command execution failed: {0}")]
    ExecutionFailed(String),
    #[error("IO error: {0}")]
    IoError(#[from] std::io::Error),
}

pub struct CommandOutput {
    pub stdout: String,
    pub stderr: String,
    pub exit_code: i32,
    pub execution_time: Duration,
}

static CONFIG: Lazy<Arc<Config>> = Lazy::new(|| {
    let args = Args::parse();
    Arc::new(load_config(&args.config).expect("Failed to load config"))
});

fn load_config(file_path: &str) -> Result<Config, Box<dyn Error>> {
    let conf = Ini::load_from_file(file_path)?;
    
    let webconf = load_section(&conf, "webconf")?;
    let paths = load_section(&conf, "paths")?;
    let commands = load_commands(&conf, &paths)?;

    Ok(Config {
        webconf,
        paths,
        commands,
    })
}

fn load_section(conf: &Ini, section_name: &str) -> Result<ConfigMap, Box<dyn Error>> {
    conf.section(Some(section_name))
        .ok_or_else(|| format!("Section '{}' not found", section_name).into())
        .map(|section| section.iter().map(|(k, v)| (k.to_string(), v.to_string())).collect())
}

fn load_commands(conf: &Ini, paths: &ConfigMap) -> Result<HashMap<String, CustomCommand>, Box<dyn Error>> {
    let commands_section = conf.section(Some("commands"))
        .ok_or_else(|| "Section 'commands' not found")?;

    commands_section.iter()
        .filter(|(key, _)| key.ends_with("Name"))
        .map(|(key, value)| {
            let base_key = key.trim_end_matches("Name");
            let command = CustomCommand {
                name: Arc::from(value.to_string()),
                value: Arc::from(commands_section.get(&format!("{}Value", base_key))
                    .ok_or_else(|| format!("{}Value not found", base_key))?.to_string()),
                command: Arc::from(replace_placeholders(
                    commands_section.get(&format!("{}Command", base_key))
                        .ok_or_else(|| format!("{}Command not found", base_key))?,
                    paths
                )),
                description: Arc::from(commands_section.get(&format!("{}Description", base_key))
                    .ok_or_else(|| format!("{}Description not found", base_key))?.to_string()),
            };
            Ok((base_key.to_string(), command)) })
        .collect()
}

fn replace_placeholders(value: &str, paths: &ConfigMap) -> String {
    let mut result = value.to_string();
    for (key, path_value) in paths {
        result = result.replace(&format!("{{{{.{}}}}}", key), path_value);
    }
    result
}

fn print_config(config: &Config) {
    println!("Web Configuration:");
    for (key, value) in &config.webconf {
        println!("  {}: {}", key, value);
    }

    println!("\nPaths:");
    for (key, value) in &config.paths {
         println!("  {}: {}", key, value);
    }

    println!("\nCommands:");
    for (name, command) in &config.commands {
        println!("{}:", name);
        println!("  Name: {}", command.name);
        println!("  Value: {}", command.value);
        println!("  Command: {}", command.command);
        println!("  Description: {}", command.description);
    }
}

async fn get_command_list(data: web::Data<Arc<Config>>) -> impl Responder {
    let api_commands: Vec<ApiCommand> = data.commands.values()
        .map(|cmd| ApiCommand {
            name: cmd.name.to_string(),
            value: cmd.value.to_string(),
            command: cmd.command.to_string(),
            description: cmd.description.to_string(),
        })
        .collect();

    HttpResponse::Ok().json(api_commands)
}

pub async fn execute_command(command: &str, timeout_duration: Duration) -> Result<CommandOutput, CommandError> {
    let start_time = Instant::now();

    info!("Executing command: {}", command);

    let result = timeout(timeout_duration, async {
        let mut cmd = if cfg!(target_os = "android") {
            StdCommand::new("/system/bin/sh")
                .arg("-c")
                .arg(command)
                .stdout(Stdio::piped())
                .stderr(Stdio::piped())
                .spawn()?
        } else if cfg!(target_os = "windows") {
            StdCommand::new("cmd")
                .args(&["/C", command])
                .stdout(Stdio::piped())
                .stderr(Stdio::piped())
                .spawn()?
        } else {
            // Assume Unix-like system (Linux, macOS, etc.)
            StdCommand::new("/bin/sh")
                .arg("-c")
                .arg(command)
                .stdout(Stdio::piped())
                .stderr(Stdio::piped())
                .spawn()?
        };

        let mut stdout = String::new();
        let mut stderr = String::new();

        if let Some(stdout_pipe) = cmd.stdout.take() {
            let reader = BufReader::new(stdout_pipe);
            for line in reader.lines() {
                let line = line?;
                info!("STDOUT: {}", line);
                stdout.push_str(&line);
                stdout.push('\n');
            }
        }

        if let Some(stderr_pipe) = cmd.stderr.take() {
            let reader = BufReader::new(stderr_pipe);
            for line in reader.lines() {
                let line = line?;
                error!("STDERR: {}", line);
                stderr.push_str(&line);
                stderr.push('\n');
            }
        }

        let status = cmd.wait()?;
        let exit_code = status.code().unwrap_or(-1);

        Ok::<_, CommandError>(CommandOutput {
            stdout,
            stderr,
            exit_code,
            execution_time: start_time.elapsed(),
        })
    }).await;

    match result {
        Ok(Ok(output)) => {
            info!("Command executed successfully in {:?}", output.execution_time);
            Ok(output)
        },
        Ok(Err(e)) => {
            error!("Command execution failed: {:?}", e);
            Err(e)
        },
        Err(_) => {
            error!("Command execution timed out after {:?}", timeout_duration);
            Err(CommandError::Timeout)
        }
    }
}

async fn run_command(
    data: web::Data<Arc<Config>>,
    path: web::Path<String>,
) -> impl Responder {
    let command_value = path.into_inner();
    match data.get_command_by_value(&command_value) {
        Some(command) => {
            match execute_command(&command.command, Duration::from_secs(30)).await {
                Ok(output) => {
                    if output.exit_code == 0 {
                        HttpResponse::Ok().json(Response {
                            success: true,
                            message: output.stdout.trim().to_string(),
                            details: None,
                        })
                    } else {
                        HttpResponse::InternalServerError().json(Response {
                            success: false, message: "Command failed".to_string(),
                            details: Some(format!("Exit code: {}, Error: {}", output.exit_code, output.stderr.trim())),
                        })
                    }
                },
                Err(e) => {
                    HttpResponse::InternalServerError().json(Response {
                        success: false,
                        message: "Command execution failed".to_string(),
                        details: Some(format!("{:?}", e)),
                    })
                },
            }
        },
        None => {
            HttpResponse::NotFound().json(Response {
                success: false,
                message: "Command not found".to_string(),
                details: Some(format!("Command with value '{}' not found", command_value)),
            })
        },
    }
}

async fn run_custom_command(
    query: web::Query<HashMap<String, String>>,
) -> impl Responder {
    if let Some(command) = query.get("custom") {
        match execute_command(command, Duration::from_secs(30)).await {
            Ok(output) => {
                if output.exit_code == 0 {
                    HttpResponse::Ok().json(Response {
                        success: true,
                        message: output.stdout.trim().to_string(),
                        details: None,
                    })
                } else {
                    HttpResponse::InternalServerError().json(Response {
                        success: false,
                        message: "Command failed".to_string(),
                        details: Some(format!("Exit code: {}, Error: {}", output.exit_code, output.stderr.trim())),
                    })
                }
            }
            Err(e) => {
                HttpResponse::InternalServerError().json(Response {
                    success: false,
                    message: "Command execution failed".to_string(),
                    details: Some(format!("{:?}", e)),
                })
            }
        }
    } else {
        HttpResponse::BadRequest().json(Response {
            success: false,
            message: "Missing 'custom' parameter".to_string(),
            details: None,
        })
    }
}

impl Config {
    fn get_command_by_value(&self, value: &str) -> Option<&CustomCommand> {
        self.commands.values().find(|cmd| cmd.value == Arc::from(value))
    }
}

// SPA handler
async fn spa_handler(path: web::Path<String>, req: HttpRequest) -> Result<HttpResponse, ActixError> {
    let path = if path.is_empty() { "index.html" } else { &path };
    
    let start = std::time::Instant::now();
    
    let response = FrontendAsset::get(path)
        .map(|content| {
            let mime = from_path(path).first_or_octet_stream();
            HttpResponse::Ok().content_type(mime.as_ref()).body(content.data.into_owned())
        })
        .unwrap_or_else(|| {
            FrontendAsset::get("index.html")
                .map(|content| HttpResponse::Ok().content_type("text/html").body(content.data.into_owned()))
                .unwrap_or_else(|| HttpResponse::NotFound().body("Not found"))
        });
    
    let duration = start.elapsed();
    info!(
        "Request duration={} method={} path={} status={}",
        duration.as_nanos(),
        req.method(),
        req.path(),
        response.status().as_u16()
    );
    
    Ok(response)
}

fn configure_logger() {
    env_logger::Builder::from_env(Env::default().default_filter_or("info"))
        .format(|buf, record| {
            let level = record.level();
            let timestamp = Local::now().format("%I:%M%p");
            writeln!(
                buf,
                "{} {} {}",
                timestamp,
                level.to_string().chars().next().unwrap(),
                record.args()
            )
        })
        .init();
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let args = Args::parse();
    
    println!("Loading configuration from: {}", args.config);
    print_config(&*CONFIG);

    configure_logger();

    HttpServer::new(move || {
        App::new()
            .wrap(Logger::new("%{r}a %s \"%{User-Agent}i\""))
            .app_data(web::Data::new(CONFIG.clone()))
            .service(
                web::scope("/api")
                    .service(web::resource("/command/list").route(web::get().to(get_command_list)))
                    .service(web::resource("/command/run/{command_value}").route(web::get().to(run_command)))
                    .service(web::resource("/command/run").route(web::get().to(run_custom_command)))
            )
            // .service(fs::Files::new("/", "./dist").index_file("index.html"))
            .service(web::resource("/{filename:.*}").route(web::get().to(spa_handler)))
    })
    .bind(format!("0.0.0.0{}", CONFIG.webconf.get("port").unwrap()))?
    .run()
    .await
}