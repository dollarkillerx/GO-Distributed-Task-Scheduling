/**
* Created by GoLand
* User: dollarkiller
* Date: 19-6-5
* Time: 下午9:06
* */
package etcd

import (
	"github.com/coreos/etcd/clientv3"
	"time"
)

var (
	Etcd *clientv3.Client
	e error
	Kv clientv3.KV
	Lease clientv3.Lease
)

func init()  {
	config := clientv3.Config{
		Endpoints:   []string{"172.17.0.2:2379"},
		DialTimeout: 5 * time.Second,
	}
	Etcd, e = clientv3.New(config)
	if e != nil {
		panic(e.Error())
	}

	Kv = clientv3.NewKV(Etcd)
	Lease = clientv3.NewLease(Etcd)

}