package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

func main() {
	users := map[string]string{
		"admin":   "1234!!",
		"AndFran": "1234!!",
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal("error connecting", err, "uri: ", os.Getenv("MONGO_URI"))
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("cannot ping", err)
	}

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	h := sha256.New()

	for user, password := range users {

		h.Write([]byte(password))
		pwd := h.Sum(nil)

		//pwd := string(h.Sum([]byte(password)))

		_, err = collection.InsertOne(ctx, bson.M{
			"username": user,
			"password": hex.EncodeToString(pwd),
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("----------")
	fmt.Println(collection)
}
