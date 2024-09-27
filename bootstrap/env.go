package bootstrap

import (
	"log"

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

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env:", err)
	}

	viper.Unmarshal(&env)
	return &env
}
