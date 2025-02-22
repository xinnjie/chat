package config

import (
	"encoding/json"
	"log"

	"github.com/spf13/viper"
)

// Credential validator config.
type ValidatorConfig struct {
	// TRUE or FALSE to set
	AddToTags bool `mapstructure:"add_to_tags"`
	//  Authentication level which triggers this validator: "auth", "anon"... or ""
	Required []string `mapstructure:"required"`
	// Validator params passed to validator unchanged.
	Config map[string]any `mapstructure:"config"`
}

// Contentx of the configuration file
type Config struct {
	// HTTP(S) address:port to listen on for websocket and long polling clients. Either a
	// numeric or a canonical name, e.g. ":80" or ":https". Could include a host name, e.g.
	// "localhost:80".
	// Could be blank: if TLS is not configured, will use ":80", otherwise ":443".
	// Can be overridden from the command line, see option --listen.
	Listen string `mapstructure:"listen"`
	// Base URL path where the streaming and large file API calls are served, default is '/'.
	// Can be overridden from the command line, see option --api_path.
	ApiPath string `mapstructure:"api_path"`
	// Cache-Control value for static content.
	CacheControl int `mapstructure:"cache_control"`
	// If true, do not attempt to negotiate websocket per message compression (RFC 7692.4).
	// It should be disabled (set to true) if you are using MSFT IIS as a reverse proxy.
	WSCompressionDisabled bool `mapstructure:"ws_compression_disabled"`
	// Address:port to listen for gRPC clients. If blank gRPC support will not be initialized.
	// Could be overridden from the command line with --grpc_listen.
	GrpcListen string `mapstructure:"grpc_listen"`
	// Enable handling of gRPC keepalives https://github.com/grpc/grpc/blob/master/doc/keepalive.md
	// This sets server's GRPC_ARG_KEEPALIVE_TIME_MS to 60 seconds instead of the default 2 hours.
	GrpcKeepalive bool `mapstructure:"grpc_keepalive_enabled"`
	// URL path for mounting the directory with static files (usually TinodeWeb).
	StaticMount string `mapstructure:"static_mount"`
	// Local path to static files. All files in this path are made accessible by HTTP.
	StaticData string `mapstructure:"static_data"`
	// Salt used in signing API keys
	APIKeySalt string `mapstructure:"api_key_salt"`
	// Maximum message size allowed from client. Intended to prevent malicious client from sending
	// very large files inband (does not affect out of band uploads).
	MaxMessageSize int `mapstructure:"max_message_size"`
	// Maximum number of group topic subscribers.
	MaxSubscriberCount int `mapstructure:"max_subscriber_count"`
	// Masked tags: tags immutable on User (mask), mutable on Topic only within the mask.
	MaskedTagNamespaces []string `mapstructure:"masked_tags"`
	// Maximum number of indexable tags.
	MaxTagCount int `mapstructure:"max_tag_count"`
	// If true, ordinary users cannot delete their accounts.
	PermanentAccounts bool `mapstructure:"permanent_accounts"`
	// URL path for exposing runtime stats. Disabled if the path is blank.
	ExpvarPath string `mapstructure:"expvar"`
	// URL path for internal server status. Disabled if the path is blank.
	ServerStatusPath string `mapstructure:"server_status"`
	// Take IP address of the client from HTTP header 'X-Forwarded-For'.
	// Useful when tinode is behind a proxy. If missing, fallback to default RemoteAddr.
	UseXForwardedFor bool `mapstructure:"use_x_forwarded_for"`
	// 2-letter country code (ISO 3166-1 alpha-2) to assign to sessions by default
	// when the country isn't specified by the client explicitly and
	// it's impossible to infer it.
	DefaultCountryCode string `mapstructure:"default_country_code"`

	// Configs for subsystems
	Cluster     ClusterConfig               `mapstructure:"cluster_config"`
	Plugins     []PluginConfig              `mapstructure:"plugins"`
	StoreConfig StoreConfig                 `mapstructure:"store_config"`
	Push        []PushConfig                `mapstructure:"push"`
	TLS         TLSConfig                   `mapstructure:"tls"`
	Auth        map[string]any              `mapstructure:"auth_config"`
	Validator   map[string]*ValidatorConfig `mapstructure:"acc_validation"`
	AccountGC   *AccountGcConfig            `mapstructure:"acc_gc_config"`
	Media       *MediaConfig                `mapstructure:"media"`
	WebRTC      CallConfig                  `mapstructure:"webrtc"`
}

func FromViper() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("/etc/tinode/config")
	viper.AddConfigPath("config/base/env")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Environment variable bindings
	envBindings := map[string]string{
		"acc_gc_config.enabled":                       "ACC_GC_ENABLED",
		"auth_config.token.key":                       "AUTH_TOKEN_KEY",
		"media.handlers.s3.access_key_id":             "AWS_ACCESS_KEY_ID",
		"media.handlers.s3.cors_origins":              "AWS_CORS_ORIGINS",
		"media.handlers.s3.region":                    "AWS_REGION",
		"media.handlers.s3.bucket":                    "AWS_S3_BUCKET",
		"media.handlers.s3.endpoint":                  "AWS_S3_ENDPOINT",
		"media.handlers.s3.secret_access_key":         "AWS_SECRET_ACCESS_KEY",
		"cluster_config.self":                         "CLUSTER_SELF",
		"acc_validation.email.config.debug_response":  "DEBUG_EMAIL_VERIFICATION_CODE",
		"default_country_code":                        "DEFAULT_COUNTRY_CODE",
		"media.handlers.fs.cors_origins":              "FS_CORS_ORIGINS",
		"webrtc.ice_servers_file":                     "ICE_SERVERS_FILE",
		"media.use_handler":                           "MEDIA_HANDLER",
		"store_config.adapters.mysql.dsn":             "MYSQL_DSN",
		"plugins.python_chat_bot.enabled":             "PLUGIN_PYTHON_CHAT_BOT_ENABLED",
		"store_config.adapters.postgres.dsn":          "POSTGRES_DSN",
		"acc_validation.email.config.auth_mechanism":  "SMTP_AUTH_MECHANISM",
		"acc_validation.email.config.domains":         "SMTP_DOMAINS",
		"acc_validation.email.config.smtp_helo_host":  "SMTP_HELO_HOST",
		"acc_validation.email.config.host_url":        "SMTP_HOST_URL",
		"acc_validation.email.config.login":           "SMTP_LOGIN",
		"acc_validation.email.config.sender_password": "SMTP_PASSWORD",
		"acc_validation.email.config.smtp_port":       "SMTP_PORT",
		"acc_validation.email.config.sender":          "SMTP_SENDER",
		"acc_validation.email.config.smtp_server":     "SMTP_SERVER",
		"store_config.use_adapter":                    "STORE_USE_ADAPTER",
		"acc_validation.tel.config.host_url":          "TEL_HOST_URL",
		"acc_validation.tel.config.sender":            "TEL_SENDER",
		"tls.autocert.email":                          "TLS_CONTACT_ADDRESS",
		"tls.enabled":                                 "TLS_ENABLED",
		"tls.autocert.domains":                        "TLS_DOMAIN_NAME",
		"push.tnpg.config.token":                      "TNPG_AUTH_TOKEN",
		"push.tnpg.config.org":                        "TNPG_ORG",
		"store_config.uid_key":                        "UID_ENCRYPTION_KEY",
		// FCM related bindings
		"push.fcm.config.enabled":          "FCM_PUSH_ENABLED",
		"push.fcm.config.project_id":       "FCM_PROJECT_ID",
		"push.fcm.config.credentials_file": "FCM_CRED_FILE",
		"push.fcm.config.api_key":          "FCM_API_KEY",
		"push.fcm.config.app_id":           "FCM_APP_ID",
		"push.fcm.config.sender_id":        "FCM_SENDER_ID",
		"push.fcm.config.vapid_key":        "FCM_VAPID_KEY",
		"push.fcm.config.android.enabled":  "FCM_INCLUDE_ANDROID_NOTIFICATION",
		"push.fcm.config.measurement_id":   "FCM_MEASUREMENT_ID",
		"webrtc.enabled":                   "WEBRTC_ENABLED",
		"server_status":                    "SERVER_STATUS_PATH",
		"push.tnpg.config.enabled":         "TNPG_PUSH_ENABLED",
	}

	for configPath, envVar := range envBindings {
		if err := viper.BindEnv(configPath, envVar); err != nil {
			log.Fatal("Cannot bind env: ", err)
		}
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Cannot load config:", err)
	}

	log.Printf("Using config file: %s", viper.ConfigFileUsed())

	return &config
}

func MustJsonRawMessage(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return json.RawMessage(b)
}
