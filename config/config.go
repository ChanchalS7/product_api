package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB        *mongo.Database
	JWTSecret string
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	JWTSecret = os.Getenv("JWT_SECRET")
}
func ConnectDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI")).SetMaxPoolSize(100).SetMinPoolSize(5).SetMaxConnIdleTime(30 * time.Minute)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ping the database to verify the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	log.Println("Successfully connected to MongoDB")
	return client

}

// var DB *mongo.Database = GetDB()

// func GetDB() *mongo.Database {
// 	client := ConnectDB()
// 	return client.Database(os.Getenv("DB_NAME"))
// }

func InitDB() {
	client := ConnectDB()
	DB = client.Database(os.Getenv("DB_NAME"))
}
