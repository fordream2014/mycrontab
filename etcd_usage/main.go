package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
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
		preValue   *mvccpb.KeyValue
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
	putResp, err = kv.Put(context.TODO(), "/mycrontab/jobs/job1", "{name: job1}", clientv3.WithPrevKV())
	if err != nil {
		fmt.Println("put failed, err: ", err)
		return
	}
	preValue = putResp.PrevKv
	fmt.Println("put success, pre value: ", string(preValue.Value))
}

func main() {
	//etcdPutUsage()
	var (
		config  clientv3.Config
		err     error
		client  *clientv3.Client
		kv      clientv3.KV
		putResp *clientv3.PutResponse
	)
	//配置
	config = clientv3.Config{
		Endpoints:   []string{"10.235.25.241:2379"},
		DialTimeout: time.Second * 5,
	}
	//连接 床见一个客户端
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	//用于读写etcd的键值对
	kv = clientv3.NewKV(client)
	putResp, err = kv.Put(context.TODO(), "/cron/jobs/job1", "bye", clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(err)
	} else {
		//获取版本信息
		fmt.Println("Revision:", putResp.Header.Revision)
		if putResp.PrevKv != nil {
			fmt.Println("key:", string(putResp.PrevKv.Key))
			fmt.Println("Value:", string(putResp.PrevKv.Value))
			fmt.Println("Version:", string(putResp.PrevKv.Version))
		}
	}
}
