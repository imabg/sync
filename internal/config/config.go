package config

import (
	"log"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Application struct {
	ErrorLog *zap.SugaredLogger
	InfoLog  *zap.SugaredLogger
	Env      Env
	Client *mongo.Client
}

type Env struct {
	ServerAddr string `mapstructure:"PORT"`
	MongoURI   string `mapstructure:"MONGO_URI"`
	DBName string `mapstucture:"DB_NAME"`
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
