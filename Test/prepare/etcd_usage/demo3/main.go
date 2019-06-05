/**
* Created by GoLand
* User: dollarkiller
* Date: 19-6-5
* Time: 下午9:06
* */
package main

import (
	"GO-Distributed-Task-Scheduling/Test/prepare/etcd_usage/demo3/etcd"
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	// lease实现锁自动过期
	// op操作
	// txn事务: if else then

	// 1.上锁(创建租约,自动续租,拿着租约去抢占一个key)
	lease, e := etcd.Lease.Grant(context.TODO(), 5)
	if e != nil {
		panic(e.Error())
	}
	leaseId := lease.ID
	// 准备一个可以用于取消自动续租的context
	ctx, cancel := context.WithCancel(context.TODO())

	leaseKeepChan, e := etcd.Lease.KeepAlive(ctx, leaseId)
	if e != nil {
		panic(e.Error())
	}

	go func() {
	forloop:
		for {
			select {
			case keepResp := <-leaseKeepChan:
				if keepResp == nil {
					fmt.Println("续租失效")
					break forloop
				}else{
					fmt.Println("收到续租答应",keepResp.ID)
				}
			}
		}
	}()

	// 抢K  (if 不存在key，then设置它，else抢锁失败)

	// 创建事务
	txn := etcd.Kv.Txn(context.TODO())
	// 定义事务
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"),"=",0)).
		Then(clientv3.OpPut("/cron/lock/job9","",clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/lock/job9")) // 否则抢锁失败
	// 提交
	txnResponse, e := txn.Commit()
	if e != nil {
		panic(e.Error())
	}
	// 判断是否抢到了锁
	if !txnResponse.Succeeded {
		fmt.Println("锁被占用",string(txnResponse.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}


	// 2.处理业务
	fmt.Println("处理任务")
	time.Sleep(5 * time.Second)


	// 3.释放锁(取消自动续租,释放租约)
	defer cancel() 	// 确保函数退出后取消自动续租
	defer etcd.Lease.Revoke(context.TODO(),leaseId) // 关闭续租
}

// 1.上锁(创建租约,自动续租,拿着租约去抢占一个key)
func lockUp() {
	lease, e := etcd.Lease.Grant(context.TODO(), 5)
	if e != nil {
		panic(e.Error())
	}
	leaseId := lease.ID
	leaseKeepChan, e := etcd.Lease.KeepAlive(context.TODO(), leaseId)
	if e != nil {
		panic(e.Error())
	}

	go func() {
		forloop:
		for {
			select {
			case keepResp := <-leaseKeepChan:
				if keepResp == nil {
					fmt.Println("续租失效")
					break forloop
				}else{
					fmt.Println("收到续租答应",keepResp.ID)
				}
			}
		}
	}()

}





