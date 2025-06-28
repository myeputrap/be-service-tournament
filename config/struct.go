package config

import (
	"time"
)

// Config is main config structure
type Config struct {
	Server            Server     `yaml:"server"`
	PostgreSQL        PostgreSQL `yaml:"postgres"`
	JWT               JWT        `yaml:"jwt"`
	DefaultLimitQuery int64      `yaml:"default_limit_query"`
	SMTP              SMTP       `yaml:"smtp"`
	HostUrl           string     `yaml:"host_url"`
	DisableEmail      bool       `yaml:"disable_email"`
}

// Server is server related config
type Server struct {
	// BasePath is router base path
	BasePath string `yaml:"base_path"`

	// LogType is log type, available value: text, json
	LogType string `yaml:"log_type"`

	// LogLevel is log level, available value: error, warning, info, debug
	LogLevel string `yaml:"log_level"`

	// HTTP is HTTP server config
	HTTP HTTPServer `yaml:"http"`

	// GRPC Port is port for grpc server
	GRPCPort string `yaml:"grpc_port"`
}

// HTTPServer is HTTP server related config
type HTTPServer struct {
	// Port is the local machine TCP Port to bind the HTTP Server to
	Port string `yaml:"port"`

	// Prefork will spawn multiple Go processes listening on the same port
	Prefork bool `yaml:"prefork"`

	// StrictRouting
	// When enabled, the router treats /foo and /foo/ as different.
	// Otherwise, the router treats /foo and /foo/ as the same.
	StrictRouting bool `yaml:"strict_routing"`

	// CaseSensitive
	// When enabled, /Foo and /foo are different routes.
	// When disabled, /Foo and /foo are treated the same.
	CaseSensitive bool `yaml:"case_sensitive"`

	// BodyLimit
	// Sets the maximum allowed size for a request body, if the size exceeds
	// the configured limit, it sends 413 - Request Entity Too Large response.
	BodyLimit int `yaml:"body_limit"`

	// Concurrency maximum number of concurrent connections
	Concurrency int `yaml:"concurrency"`

	// Timeout is HTTP server timeout
	Timeout Timeout `yaml:"timeout"`

	// AllowsOrigin
	// Put a list of origins that are allowed to access the resource,
	// separated by comma
	AllowsOrigin string `yaml:"allows_origin"`

	AssetPath string `yaml:"asset_path"`

	BaseURL string `yaml:"base_url"`

	// CacheStaticTTL defines how long (in seconds) static files should be cached.
	// This value is used to set the Cache-Control header.
	CacheStaticTTL int `yaml:"cache_static_ttl"`
}

// Timeout is server timeout related config
type Timeout struct {
	// Read is the amount of time to wait until an HTTP server
	// read operation is cancelled
	Read time.Duration `yaml:"read"`

	// Write is the amount of time to wait until an HTTP server
	// write opperation is cancelled
	Write time.Duration `yaml:"write"`

	// Read is the amount of time to wait
	// until an IDLE HTTP session is closed
	Idle time.Duration `yaml:"idle"`
}

type PostgreSQL struct {
	// Host is the PostgreSQL IP Address to connect to
	Host string `yaml:"host,omitempty"`

	// Port is the PostgreSQL Port to connect to
	Port string `yaml:"port,omitempty"`

	// Database is PostgreSQL database name
	Database string `yaml:"database"`

	// User is PostgreSQL username
	User string `yaml:"user"`

	// Password is PostgreSQL password
	Password string `yaml:"password"`

	// PathMigrate is directory for migration file
	PathMigrate string `yaml:"path_migrate"`

	// Timeone is PostgreSQL timezone
	Timezone string `yaml:"timezone"`

	// SetMaxOpenConns is maximum number of open connections to the database
	SetMaxOpenConns int `yaml:"set_max_open_conns"`

	// SetMaxIdleConns is maximum number of connections in the idle connection
	// pool
	SetMaxIdleConns int `yaml:"set_max_idle_conns"`

	// SetConnMaxIdleTime is maximum amount of time a connection may be idle
	SetConnMaxIdleTime time.Duration `yaml:"set_conn_max_idle_time"`

	// DB.SetConnMaxLifetime is maximum amount of time a connection may be
	// reused
	SetConnMaxLifetime time.Duration `yaml:"set_conn_max_lifetime"`
}

type HostPort struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type JWT struct {
	Expiration     int64  `yaml:"expiration"`
	SigningMethod  string `yaml:"signing_method"`
	SignatureKey   string `yaml:"signature_key"`
	ExpirationTemp int64  `yaml:"expiration_temp"`
}

type SMTP struct {
	Mail        string `yaml:"mail"`
	Port        int64  `yaml:"port"`
	StartTls    int64  `yaml:"start_tls"`
	TslOrSsl    int64  `yaml:"tsl_or_ssl"`
	EmailSender string `yaml:"email_sender"`
	Password    string `yaml:"password"`
}

// Default config
var defaultConfig = &Config{
	Server: Server{
		BasePath: "",
		LogType:  "json",
		LogLevel: "debug",
		HTTP: HTTPServer{
			Port:          "1360",
			Prefork:       false,
			StrictRouting: false,
			CaseSensitive: false,
			BodyLimit:     104 * 1024 * 1024,
			Concurrency:   256 * 1024,
			Timeout: Timeout{
				Read:  5,
				Write: 10,
				Idle:  120,
			},
			AllowsOrigin:   "*",
			AssetPath:      "/var/lib/padel",
			BaseURL:        "",
			CacheStaticTTL: 24 * 60 * 60,
		},
		GRPCPort: "54000",
	},
	PostgreSQL: PostgreSQL{
		Host:               "localhost",
		Port:               "15432",
		Database:           "tourney_system",
		User:               "root",
		Password:           "yourpassword",
		PathMigrate:        "file://migration",
		Timezone:           "UTC",
		SetMaxOpenConns:    0,
		SetMaxIdleConns:    2,
		SetConnMaxIdleTime: 0,
		SetConnMaxLifetime: 0,
	},
	JWT: JWT{
		Expiration:     36000000,
		SigningMethod:  "HS256",
		SignatureKey:   "kmzway87aa",
		ExpirationTemp: 600000,
	},
	DefaultLimitQuery: 100,
	SMTP: SMTP{
		Mail:        "smtp.gmail.com",
		Port:        25,
		StartTls:    587,
		TslOrSsl:    465,
		EmailSender: "",
		Password:    "",
	},
	HostUrl:      "",
	DisableEmail: true,
}
