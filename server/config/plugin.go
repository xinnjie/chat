package config

// PluginRPCFilterConfig filters for an individual RPC call. Filter strings are formatted as follows:
// <comma separated list of packet names> ; <comma separated list of topics or topic types> ; <actions (combination of C U D)>
// For instance:
// "acc,login;;CU" - grab packets {acc} or {login}; no filtering by topic, Create or Update action
// "pub,pres;me,p2p;"
type pluginRPCFilterConfig struct {
	// Filter by packet name, topic type [or exact name - not supported yet]. 2D: "pub,pres;p2p,me"
	FireHose *string `mapstructure:"fire_hose"`

	// Filter by CUD, [exact user name - not supported yet]. 1D: "C"
	Account *string `mapstructure:"account"`
	// Filter by CUD, topic type[, exact name]: "p2p;CU"
	Topic *string `mapstructure:"topic"`
	// Filter by CUD, topic type[, exact topic name, exact user name]: "CU"
	Subscription *string `mapstructure:"subscription"`
	// Filter by C.D, topic type[, exact topic name, exact user name]: "grp;CD"
	Message *string `mapstructure:"message"`

	// Call Find service, true or false
	Find bool
}

type PluginConfig struct {
	Enabled bool `mapstructure:"enabled"`
	// Unique service name
	Name string `mapstructure:"name"`
	// Microseconds to wait before timeout
	Timeout int64 `mapstructure:"timeout"`
	// Filters for RPC calls: when to call vs when to skip the call
	Filters pluginRPCFilterConfig `mapstructure:"filters"`
	// What should the server do if plugin failed: HTTP error code
	FailureCode int `mapstructure:"failure_code"`
	// HTTP Error message to go with the code
	FailureMessage string `mapstructure:"failure_text"`
	// Address of plugin server of the form "tcp://localhost:123" or "unix://path_to_socket_file"
	ServiceAddr string `mapstructure:"service_addr"`
}
