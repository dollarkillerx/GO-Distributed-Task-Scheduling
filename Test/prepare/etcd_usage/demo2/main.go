/**
* Created by GoLand
* User: dollarkiller
* Date: 19-6-3
* Time: 下午8:30
* */
package main

import (
	"GO-Distributed-Task-Scheduling/Test/prepare/etcd_usage/demo2/etcd"
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

func main() {
	//putKv()
	//getKv()
	//deleteKv()
}

// 设置kv
func putKv()  {
	//response, e := etcd.Kv.Put(context.TODO(), "/cron/jobs/hob1", "dollarkiller")
	response, e := etcd.Kv.Put(context.TODO(), "/cron/jobs/hob1", "dollarkiller",clientv3.WithPrevKV())//带历史
	if e != nil {
		panic(e.Error())
	}else{
		fmt.Println("revision:",response.Header.Revision)
		fmt.Println(string(response.PrevKv.Value)) // 返回上一个历史value 如果这次插入是第一次 就返回nil
	}
}

// 获取kv
func getKv()  {
	//response, e := etcd.Kv.Get(context.TODO(), "/cron/jobs/hob1")
	response, e := etcd.Kv.Get(context.TODO(), "/cron/jobs/", clientv3.WithPrefix())
	if e != nil {
		panic(e.Error())
	}else {
		for _,v := range response.Kvs  {
			fmt.Println(string(v.Key)," : ",string(v.Value))
		}
	}
}

// delete
func deleteKv() {
	//response, e := etcd.Kv.Delete(context.TODO(), "/cron/jobs/job2")
	//response, e := etcd.Kv.Delete(context.TODO(), "/cron/jobs/job2",clientv3.WithPrevKV())
	response, e := etcd.Kv.Delete(context.TODO(), "/cron/jobs/",clientv3.WithPrefix()) // 批量删除

	if e != nil {
		panic(e.Error())
	}
	if len(response.PrevKvs) == 0 {
		fmt.Println("zero")
	}else{
		fmt.Println(response.PrevKvs)
	}
}

