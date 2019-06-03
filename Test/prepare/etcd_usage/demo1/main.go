/**
* Created by GoLand
* User: dollarkiller
* Date: 19-6-3
* Time: 下午8:21
* */
package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
	)

	// 客户端配置
	config = clientv3.Config{
		Endpoints:[]string{"172.17.0.2:2379"},
		DialTimeout:5 * time.Second,
	}

	// 建立连接
	if client,err = clientv3.New(config);err != nil {
		panic(err.Error())
	}

	client = client
}
