package config

type tlsAutocertConfig struct {
	// Domains to support by autocert
	Domains []string `mapstructure:"domains"`
	// Name of directory where auto-certificates are cached, e.g. /etc/letsencrypt/live/your-domain-here
	CertCache string `mapstructure:"cache"`
	// Contact email for letsencrypt
	Email string `mapstructure:"email"`
}

type TLSConfig struct {
	// Flag enabling TLS
	Enabled bool `mapstructure:"enabled"`
	// Listen for connections on this address:port and redirect them to HTTPS port.
	RedirectHTTP string `mapstructure:"http_redirect"`
	// Enable Strict-Transport-Security by setting max_age > 0
	StrictMaxAge int `mapstructure:"strict_max_age"`
	// ACME autocert config, e.g. letsencrypt.org
	Autocert *tlsAutocertConfig `mapstructure:"autocert"`
	// If Autocert is not defined, provide file names of static certificate and key
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}
