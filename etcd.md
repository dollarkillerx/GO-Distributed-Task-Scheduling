### etcd小练习
- 搭建单机etcd,熟悉命令行操作
- golang调用etcd的put/get/delete/lease/watch方法
- 使用txn事务,实现分布式乐观锁

#### 搭建单机etcd,熟悉命令行操作
环境 docker centos7 llatest 
安装 https://github.com/etcd-io/etcd
`nohup ./etcd --listen-client-urls 'http://0.0.0.0:2379' --advertise-client-urls 'http://0.0.0.0:2379' &`

### 命令行简单使用:
``` 
ETCDCTL_API=3 ./etcdctl put "name" "dollarkiller"
ETCDCTL_API=3 ./etcdctl get "name"
ETCDCTL_API=3 ./etcdctl del "name"

根据前缀查询目录下的所有
ETCDCTL_API=3 ./etcdctl put "/cron/jobs/job1" "{...jb1}"
ETCDCTL_API=3 ./etcdctl put "/cron/jobs/job2" "{...jb2}"
ETCDCTL_API=3 ./etcdctl get "/cron/jobs/" --prefix

watch监听一个目录
ETCDCTL_API=3 ./etcdctl watch "/cron/jobs/" --prefix
```

### Golang操作etcd
代码:/Test/prepare/etcd_usage/demo2
`vgo get github.com/coreos/etcd/clientv3`
- 连接
``` 
package etcd

import (
	"github.com/coreos/etcd/clientv3"
	"time"
)

var (
	Etcd *clientv3.Client // 客户端
	err error
	Kv clientv3.KV  // 用于读写etcd的kv
	Lease clientv3.Lease // 租约
	Watcher clientv3.Watcher // 监听
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
	Lease = clientv3.NewLease(Etcd)
	Watcher = clientv3.NewWatcher(Etcd)
}
```
- 设置 put
``` 
func main() {
	//response, e := etcd.Kv.Put(context.TODO(), "/cron/jobs/hob1", "dollarkiller")
	response, e := etcd.Kv.Put(context.TODO(), "/cron/jobs/hob1", "dollarkiller",clientv3.WithPrevKV())//带历史
	if e != nil {
		panic(e.Error())
	}else{
		fmt.Println("revision:",response.Header.Revision)
		fmt.Println(string(response.PrevKv.Value)) // 返回上一个历史value 如果这次插入是第一次 就返回nil
	}
}
```
- get
``` 
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
```
- del
``` 
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
```
- 租约
``` 
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
```

- watch 监听
``` 
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
```

### GO操作etcd OP方式 (分布式锁基础)
代码:/Test/prepare/etcd_usage/demo3
``` 
var (
	Etcd *clientv3.Client
	e error
	Kv clientv3.KV
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

}

=======================================================================

func main() {
	putOp()
	getOp()
}

// 写
func putOp()  {
	// 创建Op:operation
	putOP := clientv3.OpPut("/cron/jobs/job8", "")
	// 执行
	opResponse, e := etcd.Kv.Do(context.TODO(), putOP)
	if e != nil {
		panic(e.Error())
	}
	fmt.Println(opResponse.Put().Header.Revision)
}

func getOp()  {
	get := clientv3.OpGet("/cron/jobs/", clientv3.WithPrefix())
	response, e := etcd.Kv.Do(context.TODO(), get)
	if e != nil {
		panic(e.Error())
	}
	fmt.Println(string(response.Get().Kvs[0].Key)," : ",string(response.Get().Kvs[0].Value))
}
```

### etcd 分布式集群乐观锁
``` 
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
```