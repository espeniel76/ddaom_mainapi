package mlogdb

import (
	"context"
	"ddaom/define"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var ctx context.Context
var cancel context.CancelFunc
var err error

const dataBase = "ddaom_log"
const col = "actions"

func RunMongodb() {
	initMongoDb()
}

func initMongoDb() {
	connect(define.DSN_MONGODB)
	if err != nil {
		fmt.Println(err)
	}
	// defer close()
}

func Close() {
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			fmt.Println(err)
		}
	}()
}

func connect(uri string) {
	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetMaxPoolSize(100)
	clientOptions.SetMinPoolSize(10)
	clientOptions.SetMaxConnIdleTime(10 * time.Second)
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	client, err = mongo.Connect(ctx, clientOptions)
}

func InsertOne(doc interface{}) (*mongo.InsertOneResult, error) {
	collection := client.Database(dataBase).Collection(col)
	result, err := collection.InsertOne(ctx, doc)
	return result, err
}

func InsertMany(docs []interface{}) (*mongo.InsertManyResult, error) {
	collection := client.Database(dataBase).Collection(col)
	result, err := collection.InsertMany(ctx, docs)
	return result, err
}

func UpdateOne(filter, update interface{}) (result *mongo.UpdateResult, err error) {
	collection := client.Database(dataBase).Collection(col)
	result, err = collection.UpdateOne(ctx, filter, update)
	return
}

func UpdateMany(filter, update interface{}) (result *mongo.UpdateResult, err error) {
	collection := client.Database(dataBase).Collection(col)
	result, err = collection.UpdateMany(ctx, filter, update)
	return
}

func DeleteOne(query interface{}) (result *mongo.DeleteResult, err error) {
	collection := client.Database(dataBase).Collection(col)
	result, err = collection.DeleteOne(ctx, query)
	return
}

func DeleteMany(query interface{}) (result *mongo.DeleteResult, err error) {
	collection := client.Database(dataBase).Collection(col)
	result, err = collection.DeleteMany(ctx, query)
	return
}
