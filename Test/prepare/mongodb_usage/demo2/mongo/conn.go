/**
* Created by GoLand
* User: dollarkiller
* Date: 19-6-7
* Time: 上午10:56
* */
package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	Mongo *mongo.Client
	MongoDb *mongo.Database
	e error
)

func init()  {
	timeout, _ := context.WithTimeout(context.Background(), 5*time.Second)
	Mongo, e = mongo.Connect(timeout, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	MongoDb = Mongo.Database("cron")
}