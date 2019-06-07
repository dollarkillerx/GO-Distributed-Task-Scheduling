/**
* Created by GoLand
* User: dollarkiller
* Date: 19-6-7
* Time: 上午10:56
* */
package main

import (
	"GO-Distributed-Task-Scheduling/Test/prepare/mongodb_usage/demo2/mongo"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// 任务执行的时间点
type TimePoint struct {
	StartTime int64 `bson:"startTime"`// 开始时间
	EndTime int64 `bson:"endTime"`// 结束时间
}

// 一条日志
type LogRecord struct {
	JobName string `bson:"jobName"` // 任务名称
	Command string `bson:"command"`// shell命令
	Err string `bson:"err"`// 脚本错误
	Content string `bson:"content"`// 脚本输出
	TimePoint TimePoint `bson:"timePoint"`// 执行时间
}

func main() {
	//insert()
	read()
}

var (
	logCollection = mongo.MongoDb.Collection("log")
)

func insert()  {
	//logCollection := mongo.MongoDb.Collection("log")

	// 插入记录
	record := &LogRecord{
		JobName:"job1",
		Command:"echo hello",
		Err:"",
		Content:"hello",
		TimePoint:TimePoint{
			StartTime:time.Now().Unix(),
			EndTime:time.Now().Unix()+10,
		},
	}

	result, e := logCollection.InsertOne(context.TODO(), record)
	if e != nil {
		panic(e.Error())
	}
	// _id:默认生成一个全局唯一id:object id:12字节的二进制
	if ids,ok := result.InsertedID.(primitive.ObjectID);ok{
		fmt.Println("ooo:",ids)
	}else{
		fmt.Println("no")
	}
}

func read()  {
	// mongodb读取回来是bson,需要反序列化

	// 按照jobName字段过滤,找出jobName=job10的数据五条
	type FindByJobName struct {
		JobName string `bson:"jobName"` // 条件
	}
	findByJobName := &FindByJobName{
		JobName: "job10",
	}

	// 查询
	cursor, e := logCollection.Find(context.TODO(), findByJobName)
	if e != nil {
		panic(e.Error())
	}
	log := make([]*LogRecord, 0)
	// 遍历结果集
	for cursor.Next(context.TODO()) {
		record := LogRecord{}
		e := cursor.Decode(&record)
		if e != nil {
			panic(e.Error())
		}
		log = append(log,&record)
	}
	fmt.Println(log)
}