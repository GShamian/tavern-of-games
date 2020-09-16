package apiserver

// Config object that store information from toml config file
type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
}

// NewConfig function. Constructor for Config
func NewConfig() *Config {
	return &Config{
		BindAddr: "*:8080",
		LogLevel: "debug",
	}
}
