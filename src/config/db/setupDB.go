package config_db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupDB() (*mongo.Database, error) {
	host := os.Getenv("DATABASE_DOMAIN")
	name := os.Getenv("DATABASE_SCHEMA")
	port := os.Getenv("DATABASE_PORT")

	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))
	if err != nil {
		panic(err)
	}

	db.Database(name).Collection("messages")
	db.Database(name).Collection("conversations")

	return db.Database(name), err
}
