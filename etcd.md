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
`vgo get github.com/coreos/etcd/clientv3`
- 连接
``` 
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
