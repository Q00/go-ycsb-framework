package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// COUNT,LOGTIME,DEVICE_NAME,ATTACK_NAME,RAW_PACKET
type logEntity struct {
	count       int    `bson:"count"`
	logtime     string `bson:"logtime"`
	attack_name string `bson:"attack_name"`
	raw_packet  string `bson:"raw_packet"`
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10)
	clientOptions := options.Client().ApplyURI("mongodb://117.17.189.6:27017").SetAuth(options.Credential{
		AuthSource: "",
		Username:   "root",
		Password:   "test123",
	})

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("몽고 DB에 연결했습니다!")

	// logsCollection := client.Database("log_db").Collection("logs")
	// insertResult, _ := logsCollection.InsertOne(context.TODO(), bson.D{
	// 	{"userID", "test1234"},
	// 	{"array", bson.A{"flying", "squirrel", "dev"}},
	// })
	// fmt.Println(insertResult)

}
