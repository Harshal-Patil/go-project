package controllers

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
)

var (
	db          *sql.DB
	redisClient *redis.Client
	ctx         context.Context
)

func Setup(database *sql.DB, redisCli *redis.Client, context context.Context) {
	db = database
	redisClient = redisCli
	ctx = context
}
func GetDB() *sql.DB {
	return db
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func GetContext() context.Context {
	return ctx
}
