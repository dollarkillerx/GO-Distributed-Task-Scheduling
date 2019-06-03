/**
* Created by GoLand
* User: dollarkiller
* Date: 19-6-3
* Time: 下午8:27
* */
package etcd

import (
	"github.com/coreos/etcd/clientv3"
	"time"
)

var (
	Etcd *clientv3.Client // 客户端
	err error
	Kv clientv3.KV  // 用于读写etcd的kv
)

func init() {
	config := clientv3.Config{
		Endpoints:[]string{"172.17.0.2:2379"},
		DialTimeout:5 * time.Second,
	}

	Etcd, err = clientv3.New(config)
	if err != nil {
		panic(err.Error())
	}
	Kv = clientv3.NewKV(Etcd)
}
