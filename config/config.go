package config

type Config struct {
	Host     string `env:"HOST"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`

	Port        string `env:"PORT"`
	LogLevel    string `env:"LOG_LEVEL"`
	MetricsPort string `env:"METRICS_PORT" envDefault:"8080"`
}
