package config

// Large file handler config.
type MediaConfig struct {
	// The name of the handler to use for file uploads.
	UseHandler string `mapstructure:"use_handler"`
	// Maximum allowed size of an uploaded file
	MaxFileUploadSize int64 `mapstructure:"max_size"`
	// Garbage collection timeout
	GcPeriod int `mapstructure:"gc_period"`
	// Number of entries to delete in one pass
	GcBlockSize int `mapstructure:"gc_block_size"`
	// Individual handler config params to pass to handlers unchanged.
	Handlers map[string]any `mapstructure:"handlers"`
}
