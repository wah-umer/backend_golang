package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	db     *mongo.Database
)

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")

	cont, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	cli, err := mongo.Connect(cont, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	if err := cli.Ping(cont, nil); err != nil {
		log.Fatalf("Mongo ping error: %v", err)
	}

	Client = cli
	db = cli.Database(dbName)

	fmt.Printf("Connected to MongoDB database %q\n", dbName)
}

func Collection(name string) *mongo.Collection {
	return db.Collection(name)
}
