package config

type clusterNodeConfig struct {
	Name string `mapstructure:"name"`
	Addr string `mapstructure:"addr"`
}

type ClusterConfig struct {
	// List of all members of the cluster, including this member
	Nodes []clusterNodeConfig `mapstructure:"nodes"`
	// Name of this cluster node
	ThisName string `mapstructure:"self"`
	// Deprecated: this field is no longer used.
	NumProxyEventGoRoutines int `mapstructure:"-"`
	// Failover configuration
	Failover ClusterFailoverConfig
}

type ClusterFailoverConfig struct {
	// Failover is enabled
	Enabled bool `mapstructure:"enabled"`
	// Time in milliseconds between heartbeats
	Heartbeat int `mapstructure:"heartbeat"`
	// Number of failed heartbeats before a leader election is initiated.
	VoteAfter int `mapstructure:"vote_after"`
	// Number of failures before a node is considered dead
	NodeFailAfter int `mapstructure:"node_fail_after"`
}
