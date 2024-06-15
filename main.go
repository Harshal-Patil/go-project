package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"go-project/controllers"
	"go-project/controllers/employee"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	db          *sql.DB
	redisClient *redis.Client
	ctx         = context.Background()
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)
	log.Println("Database connection established")

	// Initialize Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	log.Println("Redis connection established")
}

func main() {
	r := gin.Default()

	// Pass the database and Redis client to the controller
	controllers.Setup(db, redisClient, ctx)

	// Define routes
	r.POST("/upload", employee.UploadExcel)
	r.GET("/employee/:id", employee.ViewData)
	r.PUT("/employee/:id", employee.EditData)
	r.DELETE("/employee/:id", employee.DeleteData)
	r.GET("/employees", employee.ViewAllData)

	r.Run(":8080")
}
