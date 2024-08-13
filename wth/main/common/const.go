package common

const (
	WorkDir        = "/etc/wth"
	BinDir          = WorkDir + "/bin"
	BackupDir       = WorkDir + "/backup"
	TmpDir          = WorkDir + "/tmp"
	LogDir          = WorkDir + "/log"
	DomainFilePath  = WorkDir + "/domain.txt"
	SingboxDir      = WorkDir + "/sing-box"
	SingboxConfDir	= SingboxDir + "/config"
	CaddyDir        = WorkDir + "/caddy"
	CaddyFilePath   = CaddyDir + "/Caddyfile"
	CaddyServicePath = "/etc/systemd/system/caddy.service"
	SingBoxServicePath = "/etc/systemd/system/sing-box.service"
)

const CaddyServiceContent = `[Unit]
Description=Caddy
Documentation=https://caddyserver.com/docs/
After=network.target network-online.target
Requires=network-online.target

[Service]
Type=notify
User=caddy
Group=caddy
ExecStart=/usr/bin/caddy run --environ --config /etc/caddy/Caddyfile
ExecReload=/usr/bin/caddy reload --config /etc/caddy/Caddyfile --force
TimeoutStopSec=5s
LimitNOFILE=1048576
PrivateTmp=true
ProtectSystem=full
AmbientCapabilities=CAP_NET_ADMIN CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target
`

const CaddyFileContent = `# Add your Caddyfile configuration here
example.com {
    root * /usr/share/caddy
    file_server
}
`
const SingBoxServiceContent = `
[Unit]
Description=sing-box service
Documentation=https://sing-box.sagernet.org
After=network.target nss-lookup.target network-online.target

[Service]
CapabilityBoundingSet=CAP_NET_ADMIN CAP_NET_BIND_SERVICE CAP_SYS_PTRACE CAP_DAC_READ_SEARCH
AmbientCapabilities=CAP_NET_ADMIN CAP_NET_BIND_SERVICE CAP_SYS_PTRACE CAP_DAC_READ_SEARCH
ExecStart=/usr/local/bin/sing-box -D /var/lib/sing-box -C /usr/local/etc/sing-box run
ExecReload=/bin/kill -HUP $MAINPID
Restart=on-failure
RestartSec=10s
LimitNOFILE=infinity

[Install]
WantedBy=multi-user.target`
