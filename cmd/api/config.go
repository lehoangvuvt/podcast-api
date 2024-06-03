package main

import (
	"log"

	"github.com/spf13/viper"
)

func intializeDefaultConfig() {
	viper.SetDefault("PORT", 4000)
	viper.SetDefault("DB_USERNAME", "postgres")
	viper.SetDefault("DB_PASSWORD", "admin")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_NAME", "soundcom")
}

func loadEnvFile() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Cannot load .env file. Error: " + err.Error())
	}
}

func loadConfig() {
	intializeDefaultConfig()
	loadEnvFile()
}
