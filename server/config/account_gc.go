package config

// Stale unvalidated user account GC config.
type AccountGcConfig struct {
	Enabled bool `mapstructure:"enabled"`
	// How often to run GC (seconds).
	GcPeriod int `mapstructure:"gc_period"`
	// Number of accounts to delete in one pass.
	GcBlockSize int `mapstructure:"gc_block_size"`
	// Minimum hours since account was last modified.
	GcMinAccountAge int `mapstructure:"gc_min_account_age"`
}
