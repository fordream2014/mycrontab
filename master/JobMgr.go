package master

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"mycrontab/common"
	"time"
)

//任务管理器
type JobMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
}

var (
	G_JobMgr *JobMgr	//单例
)

//初始化管理器
func InitJobMgr() (err error) {
	var(
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
	)
	config = clientv3.Config{
		Endpoints: G_config.EtcdEndpoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	G_JobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}

//保存任务
func (jobMgr *JobMgr) SaveJob(job *common.Job) (old *common.Job, err error) {
	var (
		jobKey string
		jobValue []byte
		putResp *clientv3.PutResponse
		oldJobObj common.Job
	)
	//保存任务到/cron/jobs/任务名
	jobKey = common.JOB_SAVE_DIR + job.Name
	//任务信息json
	if jobValue,err = json.Marshal(job); err != nil {
		return
	}
	//保存到etcd中
	if putResp, err = jobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}
	fmt.Println("保存成功")
	//如果是更新，返回旧值
	if putResp.PrevKv != nil {
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		old = &oldJobObj
	}
	return
}
//删除任务
func (jobMgr *JobMgr) DeleteJob(name string) (old *common.Job, err error) {
	var (
		cronName string
		deleteResp *clientv3.DeleteResponse
		oldJobObj common.Job
	)
	cronName = common.JOB_SAVE_DIR + name
	if deleteResp, err = jobMgr.kv.Delete(context.TODO(), cronName, clientv3.WithPrevKV()); err != nil {
		return
	}
	if len(deleteResp.PrevKvs) != 0 {
		if err = json.Unmarshal(deleteResp.PrevKvs[0].Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		old = &oldJobObj
	}
	return
}
//获取job列表
func (jobMgr *JobMgr) GetJobList() (jobList []*common.Job, err error) {
	var (
		job *common.Job
		prefix string
		getResp *clientv3.GetResponse
		kv *mvccpb.KeyValue
	)
	prefix = common.JOB_SAVE_DIR
	if getResp, err = jobMgr.kv.Get(context.TODO(), prefix, clientv3.WithPrefix()); err != nil {
		return
	}
	jobList = make([]*common.Job, 0)
	for _,kv = range getResp.Kvs {
		job = &common.Job{}
		if err = json.Unmarshal(kv.Value, &job); err != nil {
			continue
		}
		jobList = append(jobList, job)
	}
	return
}
//杀死任务
func (jobMgr *JobMgr) JobKill(name string) (err error) {
	var (
		killName string
		leaseResp *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID
	)
	killName = common.JOB_KILL_DIR + name
	if leaseResp, err = G_JobMgr.lease.Grant(context.TODO(), 1); err != nil {
		return
	}
	leaseId = leaseResp.ID
	if _, err = G_JobMgr.kv.Put(context.TODO(), killName, "", clientv3.WithLease(leaseId)); err != nil {
		return
	}
	fmt.Println("强杀job:" + name)
	return
}






























