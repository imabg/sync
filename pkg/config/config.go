package config

import (
	"database/sql"
	"log"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	Env            Env
	MongoClient    *mongo.Database
	PostgresClient *sql.DB
}

type Env struct {
	ServerAddr     string `mapstructure:"PORT"`
	MongoURI       string `mapstructure:"MONGO_URI"`
	PostgresURI    string `mapstructure:"POSTGRES_URI"`
	DBName         string `mapstructure:"DB_NAME"`
	JwtSecretKey   string `mapstructure:"JWT_SECRET_KEY"`
	MailerHost     string `mapstructure:"MAILER_HOST"`
	MailerPort     string `mapstructure:"MAILER_PORT"`
	MailerUsername string `mapstructure:"MAILER_USERNAME"`
	MailerPassword string `mapstructure:"MAILER_PASSWORD"`
	MailerSender   string `mapstructure:"MAILER_SENDER"`
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
