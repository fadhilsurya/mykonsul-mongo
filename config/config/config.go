package config

import (
	"log"

	db "github.com/fadhilsurya/mykonsul-mongo/config/db/mongo"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	DbConfig    MongoConfig
	Db          *mongo.Database
	App         ServerConfig
	IsProd      bool
	GinMode     string
	JWTSecret   string
	RedisConfig RedisConf
}

type RedisConf struct {
	Address string
	Port    string
}

type MongoConfig struct {
	Address string
	Port    string
}

type ServerConfig struct {
	Env       string
	AppName   string
	Port      int
	YukConfig string
	RateLimit int
}

var AppConfig Config

func InitConfig() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error : %v", err)
	}

	viper.AutomaticEnv()

	AppConfig.JWTSecret = viper.GetString("JWT_SECRET")

	// get environment for database connection
	AppConfig.DbConfig.Port = viper.GetString("DB_PORT")
	AppConfig.DbConfig.Address = viper.GetString("DB_ADDRESS")

	// get environment for app environment
	AppConfig.App.Env = viper.GetString("DEV")
	AppConfig.App.Port = viper.GetInt("PORT")

	AppConfig.GinMode = viper.GetString("GIN_MODE")

	AppConfig.RedisConfig

	if AppConfig.GinMode != "release" {
		AppConfig.GinMode = ""
	}

	if AppConfig.App.Env == "PROD" {
		AppConfig.IsProd = true
	} else {
		AppConfig.IsProd = false
	}

	if AppConfig.App.Env == "" {
		AppConfig.App.Env = "dev"
	}

	if AppConfig.App.RateLimit == 0 {
		AppConfig.App.RateLimit = 5
	}

	db, err := db.ConnectMongo(AppConfig.DbConfig.Address, AppConfig.DbConfig.Port)
	if err != nil {
		log.Fatalf("error : %v", err)
		return
	}

	AppConfig.Db = db
}
