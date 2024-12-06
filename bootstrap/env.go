package bootstrap

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Env struct {
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPass         string `mapstructure:"DB_PASS"`
	DBName         string `mapstructure:"DB_NAME"`
	ContextTimeout int    `mapstructure:"CONTEXT_TIMEOUT"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
}

func NewEnv(path string) *Env {
	env := Env{}
	viper.AddConfigPath(path) // Look in the parent directory
	if os.Getenv("STK_SERVICE_RUN_MODE") == "PRODUCTION" {
		viper.SetConfigName("app.production") // The name of the file (without extension)
	} else {
		viper.SetConfigName("app.develop") // The name of the file (without extension)
	}

	viper.SetConfigType("env") // Set the file type to ".env"
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env:", err)
	}

	viper.Unmarshal(&env)
	return &env
}
