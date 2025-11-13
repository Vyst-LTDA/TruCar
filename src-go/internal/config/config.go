package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
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
			// Config file not found; ignore error if desired
			log.Println("Could not find .env file, using default values")
		} else {
			log.Fatalf("Error reading config file: %s", err)
		}
	}

	// Create a dummy .env file if it doesn't exist
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		f, err := os.Create(".env")
		if err != nil {
			log.Fatalf("Error creating .env file: %s", err)
		}
		defer f.Close()
		// Write default values to .env file
		f.WriteString("DB_DSN=test.db\n")
		f.WriteString("JWT_SECRET=your-secret-key\n")
		f.WriteString("SERVER_PORT=8080\n")
	}

	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}
