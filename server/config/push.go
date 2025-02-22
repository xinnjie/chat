package config

type PushConfig struct {
	Name   string `mapstructure:"name"`
	Config any    `mapstructure:"config"`
}
