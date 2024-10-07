use std::fs::{self, FileType};
use std::io;
use std::os::unix::fs::MetadataExt;
use std::path::{Path, PathBuf};
use std::time::SystemTime;
use chrono::{DateTime, Local};
use users::{get_user_by_uid, get_group_by_gid};
use serde::{Serialize, Deserialize};
use serde_json::to_string_pretty;
use rayon::prelude::*;
use thiserror::Error;

mod models {
    use super::*;

    #[derive(Serialize, Deserialize, Clone)]
    pub struct FileInfo {
        pub name: String,
        pub path: String,
        pub is_dir: bool,
        pub size: u64,
        pub last_modified: String,
        pub is_editable: bool,
        pub permissions: String,
        pub file_type: String,
        pub owner: String,
        pub group: String,
        pub creation_date: String,
    }

    #[derive(Serialize, Deserialize)]
    pub struct DirectoryListing {
        pub current_path: String,
        pub previous_path: String,
        pub files: Vec<FileInfo>,
    }
}

use models::{FileInfo, DirectoryListing};

#[derive(Error, Debug)]
pub enum ListingError {
    #[error("IO error: {0}")]
    Io(#[from] io::Error),
    #[error("Serialization error: {0}")]
    Serialization(#[from] serde_json::Error),
}

pub struct ListOptions {
    storage_dir: PathBuf,
    show_hidden: bool,
}

impl ListOptions {
    pub fn new(storage_dir: PathBuf) -> Self {
        ListOptions {
            storage_dir,
            show_hidden: false,
        }
    }

    pub fn show_hidden(mut self, show: bool) -> Self {
        self.show_hidden = show;
        self
    }
}

pub trait FileInfoProvider {
    fn get_file_info(&self, path: &Path) -> Result<FileInfo, ListingError>;
}

pub struct DefaultFileInfoProvider;

impl FileInfoProvider for DefaultFileInfoProvider {
    fn get_file_info(&self, path: &Path) -> Result<FileInfo, ListingError> {
        let metadata = fs::metadata(path)?;
        let file_name = path.file_name().unwrap_or_default().to_string_lossy().into_owned();
        
        let owner = get_user_by_uid(metadata.uid())
            .map(|user| user.name().to_string_lossy().into_owned())
            .unwrap_or_else(|| metadata.uid().to_string());
        
        let group = get_group_by_gid(metadata.gid())
            .map(|group| group.name().to_string_lossy().into_owned())
            .unwrap_or_else(|| metadata.gid().to_string());

        Ok(FileInfo {
            name: file_name,
            path: path.to_string_lossy().into_owned(),
            is_dir: metadata.is_dir(),
            size: if metadata.is_dir() { 4096 } else { metadata.len() },
            last_modified: format_time(metadata.modified()?),
            is_editable: metadata.mode() & 0o200 != 0,
                        permissions: format!("{:04o}", metadata.mode() & 0o777),
            file_type: get_file_type(&metadata.file_type()),
            owner,
            group,
            creation_date: format_time(metadata.created()?),
        })
    }
}

fn get_file_type(file_type: &FileType) -> String {
    if file_type.is_dir() {
        "directory".to_string()
    } else if file_type.is_file() {
        "file".to_string()
    } else if file_type.is_symlink() {
        "symlink".to_string()
    } else {
        "unknown".to_string()
    }
}

fn format_time(time: SystemTime) -> String {
    let datetime: DateTime<Local> = time.into();
    datetime.format("%Y-%m-%d %H:%M:%S").to_string()
}

pub fn list_directory(options: &ListOptions, provider: &impl FileInfoProvider) -> Result<DirectoryListing, ListingError> {
    let files: Vec<FileInfo> = fs::read_dir(&options.storage_dir)?
        .into_iter()
        .filter_map(|entry| entry.ok())
        .filter(|entry| {
            options.show_hidden || !entry.file_name().to_string_lossy().starts_with('.')
        })
        .filter_map(|entry| provider.get_file_info(&entry.path()).ok())
        .filter(|file_info| file_info.is_dir)
        .collect();

    let mut sorted_files = files;
    sorted_files.par_sort_by(|a, b| a.name.cmp(&b.name));

    let current_path = options.storage_dir.to_string_lossy().into_owned();
    let previous_path = options.storage_dir.parent()
        .map(|p| p.to_string_lossy().into_owned())
        .unwrap_or_else(|| "/".to_string());

    Ok(DirectoryListing {
        current_path,
        previous_path,
        files: sorted_files,
    })
}

pub fn generate_json_output(listing: &DirectoryListing) -> Result<String, ListingError> {
    Ok(to_string_pretty(listing)?)
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::HashMap;

    struct MockFileInfoProvider {
        files: HashMap<String, FileInfo>,
    }

    impl FileInfoProvider for MockFileInfoProvider {
        fn get_file_info(&self, path: &Path) -> Result<FileInfo, ListingError> {
            self.files.get(path.to_str().unwrap())
                .cloned()
                .ok_or_else(|| ListingError::Io(io::Error::new(io::ErrorKind::NotFound, "File not found")))
        }
    }

    #[test]
    fn test_list_directory() {
        let mut mock_files = HashMap::new();
        mock_files.insert("/home/user/dir1".to_string(), FileInfo {
            name: "dir1".to_string(),
            path: "/home/user/dir1".to_string(),
            is_dir: true,
            size: 4096,
            last_modified: "2023-01-01 00:00:00".to_string(),
            is_editable: true,
            permissions: "0755".to_string(),
            file_type: "directory".to_string(),
            owner: "user".to_string(),
            group: "user".to_string(),
            creation_date: "2023-01-01 00:00:00".to_string(),
        });

        let provider = MockFileInfoProvider { files: mock_files };
        let options = ListOptions::new(PathBuf::from("/home/user"));
        let result = list_directory(&options, &provider);

        assert!(result.is_ok());
                let listing = result.unwrap();
        assert_eq!(listing.current_path, "/home/user");
        assert_eq!(listing.previous_path, "/home");
        assert_eq!(listing.files.len(), 1);
        assert_eq!(listing.files[0].name, "dir1");
    }
}

fn main() -> Result<(), ListingError> {
    let options = ListOptions::new(PathBuf::from("/home/pew/wth"))
        .show_hidden(false);

    let provider = DefaultFileInfoProvider;
    let directory_listing = list_directory(&options, &provider)?;

    let json_output = generate_json_output(&directory_listing)?;
    println!("{}", json_output);

    Ok(())
}