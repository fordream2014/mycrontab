package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"

	//"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

/*
使用etcd的put功能
*/
func etcdPutUsage() {
	var (
		Config     clientv3.Config
		EtcdClient *clientv3.Client
		err        error
		putResp    *clientv3.PutResponse
		kv         clientv3.KV
	)

	Config = clientv3.Config{
		Endpoints:   []string{"10.235.25.241:2379"},
		DialTimeout: 5 * time.Second,
	}
	if EtcdClient, err = clientv3.New(Config); err != nil {
		fmt.Println("connect fail, err: ", err)
		return
	}
	fmt.Println("connect success!")
	defer EtcdClient.Close()

	kv = clientv3.NewKV(EtcdClient)
	putResp, err = kv.Put(context.TODO(), "/mycrontab/jobs/job2", "{name: job2}", clientv3.WithPrevKV())
	if err != nil {
		fmt.Println("put failed, err: ", err)
		return
	}
	fmt.Println("Revision:", putResp.Header.Revision)
	fmt.Println("Value:", string(putResp.PrevKv.Value))
}

func etcdGetUsage() {
	var (
		Config     clientv3.Config
		EtcdClient *clientv3.Client
		err        error
		getResp    *clientv3.GetResponse
		kv         clientv3.KV
	)

	Config = clientv3.Config{
		Endpoints:   []string{"10.235.25.241:2379"},
		DialTimeout: 5 * time.Second,
	}
	if EtcdClient, err = clientv3.New(Config); err != nil {
		fmt.Println("connect fail, err: ", err)
		return
	}
	fmt.Println("connect success!")
	defer EtcdClient.Close()

	kv = clientv3.NewKV(EtcdClient)
	getResp, err = kv.Get(context.TODO(), "/mycrontab/jobs/job1")
	if err != nil {
		fmt.Println("get failed, err: ", err)
		return
	}
	fmt.Println(getResp.Kvs)
}

func etcdGetCountUsage() {
	var (
		Config     clientv3.Config
		EtcdClient *clientv3.Client
		err        error
		getResp    *clientv3.GetResponse
		kv         clientv3.KV
	)

	Config = clientv3.Config{
		Endpoints:   []string{"10.235.25.241:2379"},
		DialTimeout: 5 * time.Second,
	}
	if EtcdClient, err = clientv3.New(Config); err != nil {
		fmt.Println("connect fail, err: ", err)
		return
	}
	fmt.Println("connect success!")
	defer EtcdClient.Close()

	kv = clientv3.NewKV(EtcdClient)
	getResp, err = kv.Get(context.TODO(), "/mycrontab/jobs/", clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		fmt.Println("get count failed, err: ", err)
		return
	}
	fmt.Println(getResp.Kvs, getResp.Count)
}

//租约
func etcdLease() {
	var (
		Config     clientv3.Config
		EtcdClient *clientv3.Client
		err        error
		getResp    *clientv3.GetResponse
		putResp		*clientv3.PutResponse
		kv         clientv3.KV
		lease 	   clientv3.Lease
		leaseResp  *clientv3.LeaseGrantResponse
		leaseid		clientv3.LeaseID
	)

	Config = clientv3.Config{
		Endpoints:   []string{"10.235.25.241:2379"},
		DialTimeout: 5 * time.Second,
	}
	if EtcdClient, err = clientv3.New(Config); err != nil {
		fmt.Println("connect fail, err: ", err)
		return
	}
	fmt.Println("connect success!")
	defer EtcdClient.Close()

	//申请一个lease租约
	lease = clientv3.NewLease(EtcdClient)
	//申请一个10秒的租约
	if leaseResp,err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}
	leaseid = leaseResp.ID

	kv = clientv3.NewKV(EtcdClient)
	if putResp, err = kv.Put(context.TODO(), "/mycrontab/jobs/job3", "job3", clientv3.WithLease(leaseid)); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("写入成功", putResp.Header.Revision)

	for {
		getResp, err = kv.Get(context.TODO(), "/mycrontab/jobs/job3")
		if err != nil {
			fmt.Println("get count failed, err: ", err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("kv过期了")
			break
		}
		fmt.Println("还没有过期:", getResp.Kvs)
		time.Sleep(2*time.Second)
	}
}

//租约，自动续租
func etcdKeepAliveLease() {
	var (
		Config     clientv3.Config
		EtcdClient *clientv3.Client
		err        error
		getResp    *clientv3.GetResponse
		putResp		*clientv3.PutResponse
		kv         clientv3.KV
		lease 	   clientv3.Lease
		leaseResp  *clientv3.LeaseGrantResponse
		leaseid		clientv3.LeaseID
		leaseChan	<-chan *clientv3.LeaseKeepAliveResponse
		keepresp *clientv3.LeaseKeepAliveResponse
	)

	Config = clientv3.Config{
		Endpoints:   []string{"10.235.25.241:2379"},
		DialTimeout: 5 * time.Second,
	}
	if EtcdClient, err = clientv3.New(Config); err != nil {
		fmt.Println("connect fail, err: ", err)
		return
	}
	fmt.Println("connect success!")
	defer EtcdClient.Close()

	//申请一个lease租约
	lease = clientv3.NewLease(EtcdClient)
	//申请一个10秒的租约
	if leaseResp,err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}
	leaseid = leaseResp.ID

	//自动续租
	leaseChan,err = lease.KeepAlive(context.TODO(), leaseid)
	go func() {
		for {
			select {
			case keepresp = <- leaseChan :
				if keepresp == nil {
					fmt.Println("租约以失效")
					goto END
				} else {
					fmt.Println("收到自动续租的应答")
				}
			}
		}
		END:
	}()

	kv = clientv3.NewKV(EtcdClient)
	if putResp, err = kv.Put(context.TODO(), "/mycrontab/jobs/job4", "job3", clientv3.WithLease(leaseid)); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("写入成功", putResp.Header.Revision)

	for {
		getResp, err = kv.Get(context.TODO(), "/mycrontab/jobs/job4")
		if err != nil {
			fmt.Println("get count failed, err: ", err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("kv过期了")
			break
		}
		fmt.Println("还没有过期:", getResp.Kvs)
		time.Sleep(2*time.Second)
	}
}

//watch使用
func etcdWatchUsage() {
	var (
		Config     clientv3.Config
		EtcdClient *clientv3.Client
		err        error
		getResp    *clientv3.GetResponse
		kv         clientv3.KV
		revision	int64
		watcher     clientv3.Watcher
		watchRespChan  clientv3.WatchChan
		watchResp		clientv3.WatchResponse
		ctx			context.Context
		cancelFunc 	context.CancelFunc
		event 		*clientv3.Event
	)

	Config = clientv3.Config{
		Endpoints:   []string{"10.235.25.241:2379"},
		DialTimeout: 5 * time.Second,
	}
	if EtcdClient, err = clientv3.New(Config); err != nil {
		fmt.Println("connect fail, err: ", err)
		return
	}
	fmt.Println("connect success!")
	defer EtcdClient.Close()

	kv = clientv3.NewKV(EtcdClient)
	if getResp, err = kv.Get(context.TODO(), "/mycrontab/jobs/job4"); err != nil {
		fmt.Println(err)
		return
	}
	//当前etcd集群事务ID，单调递增
	revision = getResp.Header.Revision + 1
	//创建一个watcher
	watcher = clientv3.NewWatcher(EtcdClient)
	ctx, cancelFunc = context.WithCancel(context.TODO())
	time.AfterFunc(10 * time.Second, func() {
		cancelFunc()
	})

	watchRespChan = watcher.Watch(ctx, "/mycrontab/jobs/job4", clientv3.WithRev(revision))
	for watchResp = range watchRespChan {
		for _,event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT :
				fmt.Println("修改为：", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE :
				fmt.Println("删除了", "Revision:", event.Kv.ModRevision)
			}
		}
	}
}

//分布式锁
func distributeLock() {
	var (
		Config     clientv3.Config
		EtcdClient *clientv3.Client
		err        error
		ctx	context.Context
		cancelFunc 	context.CancelFunc
		lease	clientv3.Lease
		leaseGrantResp	*clientv3.LeaseGrantResponse
		leaseKeepAliveChan  <-chan *clientv3.LeaseKeepAliveResponse
		leaseKeepAliveResp *clientv3.LeaseKeepAliveResponse
		kv clientv3.KV
		txn clientv3.Txn
		txnResp *clientv3.TxnResponse
		leaseId clientv3.LeaseID
	)

	Config = clientv3.Config{
		Endpoints:   []string{"10.235.25.241:2379"},
		DialTimeout: 5 * time.Second,
	}
	if EtcdClient, err = clientv3.New(Config); err != nil {
		fmt.Println("connect fail, err: ", err)
		return
	}
	fmt.Println("连接etcd成功")
	defer EtcdClient.Close()

	//创建租约
	lease = clientv3.NewLease(EtcdClient)
	//申请5秒的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}
	leaseId = leaseGrantResp.ID
	fmt.Println(leaseId)

	//准备一个用于取消自动续租的context
	ctx, cancelFunc = context.WithCancel(context.TODO())
	//确保函数退出后，自动续租会停止
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)

	//自动续约
	if leaseKeepAliveChan, err = lease.KeepAlive(ctx, leaseGrantResp.ID); err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		for {
			select {
			case leaseKeepAliveResp = <-leaseKeepAliveChan :
				if leaseKeepAliveChan != nil {
					fmt.Println("续租成功")
				} else {
					fmt.Println("续租失效")
					goto END
				}
			}
			time.Sleep(1*time.Second)
		}

		END:
	}()

	kv = clientv3.NewKV(EtcdClient)
	//创建事务
	txn = kv.Txn(context.TODO())
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"), "=", 0)).
		Then(clientv3.OpPut("/cron/lock/job9", "", clientv3.WithLease(leaseGrantResp.ID))).
		Else(clientv3.OpGet("/cron/lock/job9"))

	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}

	if !txnResp.Succeeded {
		fmt.Println("锁被占用", txnResp.Responses[0].GetResponseRange().Kvs[0].Value)
		return
	}
	//处理
	fmt.Println("执行任务")
	time.Sleep(5 * time.Second)

	//释放锁
	//defer会把租约释放掉，关联的kv就被删除了
}

func main() {
	distributeLock()
	//etcdWatchUsage()
	//etcdKeepAliveLease()
	//etcdLease()
	//etcdGetCountUsage()
	//etcdPutUsage()
	//etcdGetUsage()
}
