package config

import (
	"log"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Application struct {
	ErrorLog *zap.SugaredLogger
	InfoLog  *zap.SugaredLogger
	Env      Env
}

type Env struct {
	ServerAddr string `mapstructure:"PORT"`
	MongoURI   string `mapstructure:"MONGO_URI"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Can't file .env file: %v", err)
	}
	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatalf("Can't load env: %v", err)
	}
	return &env
}
