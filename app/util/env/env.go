package env

import (
	"log"

	"github.com/spf13/viper"
)

func InitEnv() {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../..") // Adjust path if needed

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading env file: %v", err)
	}

	viper.AutomaticEnv()
}

// Example getter functions
func GetDBHost() string     { return viper.GetString("DB_HOST") }
func GetDBPort() string     { return viper.GetString("DB_PORT") }
func GetDBUser() string     { return viper.GetString("DB_USER") }
func GetDBPassword() string { return viper.GetString("DB_PASSWORD") }
func GetDBName() string     { return viper.GetString("DB_NAME") }
func GetAPIKey() string     { return viper.GetString("API_KEY") }
