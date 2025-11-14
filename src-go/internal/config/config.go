package config

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go-api/internal/logging"
)

type Config struct {
	DB_DSN        string `mapstructure:"DB_DSN"`
	JWT_SECRET    string `mapstructure:"JWT_SECRET"`
	SERVER_PORT   string `mapstructure:"SERVER_PORT"`
	REDIS_ADDR    string `mapstructure:"REDIS_ADDR"`
	REDIS_PASSWORD string `mapstructure:"REDIS_PASSWORD"`
	REDIS_DB      int    `mapstructure:"REDIS_DB"`
}

var AppConfig *Config

func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("DB_DSN", "test.db")
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logging.Logger.Warn("Could not find .env file, using default values")
		} else {
			logging.Logger.Fatal("Error reading config file", zap.Error(err))
		}
	}

	// Create a dummy .env file if it doesn't exist
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		f, err := os.Create(".env")
		if err != nil {
			logging.Logger.Fatal("Error creating .env file", zap.Error(err))
		}
		defer f.Close()
		// Write default values to .env file
		f.WriteString("DB_DSN=test.db\n")
		f.WriteString("JWT_SECRET=your-secret-key\n")
		f.WriteString("SERVER_PORT=8080\n")
	}

	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		logging.Logger.Fatal("Unable to decode into struct", zap.Error(err))
	}
}
