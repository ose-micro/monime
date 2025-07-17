package monime

type Config struct {
	BaseURL    string `mapstructure:"base_url"`
	Access     string `mapstructure:"access"`
	Space      string `mapstructure:"space"`
	Version    string `mapstructure:"version"`
	TimeoutSec int    `mapstructure:"timeout_sec"`
}
