package monime

type Config struct {
	BaseURL    string `mapstructure:"base_url"`
	Access     string `mapstructure:"access"`
	Space      string `mapstructure:"space"`
	TimeoutSec int    `mapstructure:"timeout_sec"`
}
