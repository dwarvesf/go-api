package config

import (
	"github.com/spf13/viper"
)

// Loader load config from reader into Viper
type Loader interface {
	Load(viper.Viper) (*viper.Viper, error)
}

type Config struct {
	App            string
	Env            string
	Version        string
	ServerName     string
	BaseURL        string
	Port           string
	AllowedOrigins string
	SecretKey      string

	// log system
	SentryDSN string
}

func (c *Config) IsLocal() bool {
	return c.Env == "local"
}

type ENV interface {
	GetBool(string) bool
	GetString(string) string
}

func Generate(v ENV) *Config {
	return &Config{
		App:            v.GetString("APP"),
		Env:            v.GetString("ENV"),
		SecretKey:      v.GetString("SECRET_KEY"),
		SentryDSN:      v.GetString("SENTRY_DSN"),
		Version:        v.GetString("VERSION"),
		ServerName:     v.GetString("SERVER_NAME"),
		BaseURL:        v.GetString("BASE_URL"),
		Port:           v.GetString("PORT"),
		AllowedOrigins: v.GetString("ALLOWED_ORIGINS"),
	}
}

func DefaultConfigLoaders() []Loader {
	var loaders []Loader
	fileLoader := NewDotenv()
	loaders = append(loaders, fileLoader)

	return loaders
}

// LoadConfig load config from loader list
func LoadConfig(loaders []Loader) *Config {
	v := viper.New()
	v.SetDefault("APP", "scc")
	v.SetDefault("PORT", "3000")
	v.SetDefault("ENV", "local")
	v.SetDefault("ALLOWED_ORIGINS", "*")
	v.SetDefault("VERSION", "0.0.1")
	v.SetDefault("SERVER_NAME", "local")

	for idx := range loaders {
		newV, err := loaders[idx].Load(*v)

		if err == nil {
			v = newV
		}
	}
	return Generate(v)
}

func LoadTestConfig() Config {
	return Config{
		Port: "8080",
	}
}
