package config

type StoreConfig struct {
	// 16-byte key for XTEA. Used to initialize types.UidGenerator.
	UidKey string `mapstructure:"uid_key"`
	// Maximum number of results to return from adapter.
	MaxResults int `mapstructure:"max_results"`
	// DB adapter name to use. Should be one of those specified in `Adapters`.
	UseAdapter string `mapstructure:"use_adapter"`
	// Configurations for individual adapters.
	Adapters map[string]any `mapstructure:"adapters"`
}
