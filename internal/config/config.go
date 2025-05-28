package config

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"log"
	tweetsModels "twitter-api/internal/tweets/models"
	usersModels "twitter-api/internal/users/models"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type Config struct {
	Port         string
	Env          string
	Postgresdsn  string
	RedisOptions *redis.Options
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("[Config] Error reading config: %v", err)
	}

	return &Config{
		Port:        viper.GetString("server.port"),
		Env:         viper.GetString("env"),
		Postgresdsn: viper.GetString("db.postgres.dsn"),
		RedisOptions: &redis.Options{
			Addr:     viper.GetString("db.redis.addr"),
			Password: viper.GetString("db.redis.password"),
			DB:       viper.GetInt("db.redis.db"),
		},
	}
}

func (c *Config) Postgres() *gorm.DB {
	db, err := gorm.Open(postgres.Open(c.Postgresdsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[Config] Error connecting to postgres: %v", err)
	}

	if err := db.AutoMigrate(&usersModels.User{}, &usersModels.Follower{}, &tweetsModels.Tweet{}); err != nil {
		log.Fatalf("[Config] Error migrating tables: %v", err)
	}
	db.Exec("PRAGMA foreign_keys = ON;")

	return db
}

func (c *Config) Redis() *redis.Client {
	rdb := redis.NewClient(c.RedisOptions)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("[Config] Error connecting to redis: %v", err)
	}

	return rdb
}
