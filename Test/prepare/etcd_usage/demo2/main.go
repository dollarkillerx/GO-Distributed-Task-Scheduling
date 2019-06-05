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
	"github.com/coreos/etcd/mvcc/mvccpb"
	"strconv"
	"time"
)

func main() {
	//putKv()
	//getKv()
	//deleteKv()
	//leaseGrant()
	watch()
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

// 租约
func leaseGrant()  {
	les, e := etcd.Lease.Grant(context.TODO(), 10) // ttl 秒
	if e != nil {
		panic(e.Error())
	}

	// Put 一个KV 与租约关联实现10秒后过期
	leaseId := les.ID // 获取租约的id
	putResponse, e := etcd.Kv.Put(context.TODO(), "/cron/lock/job1", "", clientv3.WithLease(leaseId))
	if e != nil {
		panic(e.Error())
	}

	// 自动续租
	etcdCh, e := etcd.Lease.KeepAlive(context.TODO(), leaseId)
	if e != nil {
		panic(e.Error())
	}

	// 处理续租应答的协程
	go func() {
		forloop:
		for {
			select {
			case keepResp := <- etcdCh:
				if keepResp == nil {
					fmt.Println("租约已经失效")
					break forloop
				}else{ // 每秒租约一次,所以就会受到一次应答
					fmt.Println("收到自动租约应答:",keepResp.ID)
				}
 			}
		}
	}()


	fmt.Println("写入成功",putResponse.Header.Revision)

	// 定时看一下kv过期没有
	for {
		getResponse, e := etcd.Kv.Get(context.TODO(), "/cron/lock", clientv3.WithPrefix())
		if e != nil {
			panic(e.Error())
		}

		if getResponse.Count == 0 {
			fmt.Println("Kv 过期")
			break
		}

		for _,v := range getResponse.Kvs {
			fmt.Println("还没有过期:",string(v.Key) ," : ",string(v.Value))
		}
		time.Sleep(time.Millisecond * 500)
	}
}

// 监听
func watch() {
	// 模拟etcd中KV的变换
	go func() {
		i := 0
		for  {
			i++
			etcd.Kv.Put(context.TODO(),"/cron/jobs/job7","这个是第"+strconv.Itoa(i))
			etcd.Kv.Delete(context.TODO(),"/cron/jobs/job7")
			time.Sleep(time.Second)
		}
	}()

	// 向get到当前的value,并监听后续的变换
	response, e := etcd.Kv.Get(context.TODO(), "/cron/jobs/job7")
	if e != nil {
		panic(e.Error())
	}
	if response.Count != 0 {
		fmt.Println(string(response.Kvs[0].Key)," : ",string(response.Kvs[0].Value))
	}
	// 当前etcd集群事务id,单调递增
	version := response.Header.Revision + 1

	// 创建一个watcher
	chans := etcd.Watcher.Watch(context.TODO(), "/cron/jobs/job7", clientv3.WithRev(version)) // 第三个参数为监听的起点版本
	go func() {
		for {
			select {
			case watchResponse := <-chans:
				for _,event := range watchResponse.Events {
					switch event.Type {
					case mvccpb.PUT:
						fmt.Println("Put 修改",string(event.Kv.Value),"Revision:",event.Kv.CreateRevision," Mod: ",event.Kv.ModRevision)
					case mvccpb.DELETE:
						fmt.Println("Delete 删除",string(event.Kv.Value),"Revision:",event.Kv.ModRevision)
					}
				}
			}
		}
	}()


	time.Sleep(10 * time.Second)
}