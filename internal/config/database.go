package config

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLConnection(cfg Config) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		time.Sleep(2 * time.Second)
	}
	return nil, err
}

func NewRedisClient(cfg Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	var err error
	for i := 0; i < 10; i++ {
		_, err = client.Ping(context.Background()).Result()
		if err == nil {
			return client, nil
		}
		time.Sleep(2 * time.Second)
	}

	return nil, err
}
