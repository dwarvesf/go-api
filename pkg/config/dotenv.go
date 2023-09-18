package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// FileReader read config from file
type dotenv struct {
}

// NewDotenv create new file loader with filename and dirname
func NewDotenv() Loader {
	return &dotenv{}
}

// Load from yml file
func (r *dotenv) Load(v viper.Viper) (*viper.Viper, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &v, nil
}
