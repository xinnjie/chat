package config

type CallConfig struct {
	// Enable video/voice calls.
	Enabled bool `mapstructure:"enabled"`
	// Timeout in seconds before a call is dropped if not answered.
	CallEstablishmentTimeout int `mapstructure:"call_establishment_timeout"`
	// ICE servers.
	ICEServers []IceServer `mapstructure:"ice_servers"`
	// Alternative config as an external file.
	ICEServersFile string `mapstructure:"ice_servers_file"`
}

// ICE server config.
type IceServer struct {
	Username       string   `mapstructure:"username,omitempty"`
	Credential     string   `mapstructure:"credential,omitempty"`
	CredentialType string   `mapstructure:"credential_type,omitempty"`
	Urls           []string `mapstructure:"urls,omitempty"`
}
