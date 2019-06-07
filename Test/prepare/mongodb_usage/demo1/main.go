/**
* Created by GoLand
* User: dollarkiller
* Date: 19-6-7
* Time: 上午10:32
* */
package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func main() {
	// 1. 建立连接
	timeout, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, e := mongo.Connect(timeout, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if e != nil {
		panic(e.Error())
	}
	// 2. 选择数据库my_db
	database := client.Database("my_db")
	collection := database.Collection("my_collection")
	collection = collection
	// 3. 选择表my_collection

}
