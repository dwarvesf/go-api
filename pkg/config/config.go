package config

import (
	"github.com/spf13/viper"
)

const (
	serverVersion = "0.0.1"
)

// Loader load config from reader into Viper
type Loader interface {
	Load(viper.Viper) (*viper.Viper, error)
}

// Config struct for config
type Config struct {
	App            string
	Env            string
	Version        string
	ServerName     string
	BaseURL        string
	Port           string
	AllowedOrigins string
	SecretKey      string
	DatabaseURL    string
	DBMaxOpenConns int
	DBMaxIdleConns int

	// log system
	SentryDSN string
}

// IsLocal check if env is local
func (c *Config) IsLocal() bool {
	return c.Env == "local"
}

// ENV interface for getting env
type ENV interface {
	GetBool(string) bool
	GetString(string) string
	GetInt(string) int
}

// Generate generate config from ENV
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
		DatabaseURL:    v.GetString("DATABASE_URL"),
		DBMaxOpenConns: v.GetInt("DB_MAX_OPEN_CONNS"),
		DBMaxIdleConns: v.GetInt("DB_MAX_IDLE_CONNS"),
	}
}

// DefaultConfigLoaders return default config loaders
func DefaultConfigLoaders() []Loader {
	var loaders []Loader
	fileLoader := NewDotenv()
	loaders = append(loaders, fileLoader)

	return loaders
}

// LoadConfig load config from loader list
func LoadConfig(loaders []Loader) *Config {
	v := viper.New()
	v.SetDefault("APP", "go-api")
	v.SetDefault("PORT", "3000")
	v.SetDefault("ENV", "prod")
	v.SetDefault("ALLOWED_ORIGINS", "*")
	v.SetDefault("VERSION", serverVersion)
	v.SetDefault("SERVER_NAME", "local")
	v.SetDefault("DB_MAX_OPEN_CONNS", 10)
	v.SetDefault("DB_MAX_IDLE_CONNS", 5)

	for idx := range loaders {
		newV, err := loaders[idx].Load(*v)

		if err == nil {
			v = newV
		}
	}

	v.AutomaticEnv()

	return Generate(v)
}

// LoadTestConfig load test config
func LoadTestConfig() Config {
	return Config{
		App:         "go-api",
		Env:         "test",
		Version:     serverVersion,
		DatabaseURL: "postgres://postgres:postgres@localhost:5433/go-api-db-test?sslmode=disable",
		Port:        "4000",
	}
}
