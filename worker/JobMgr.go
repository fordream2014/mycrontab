package worker

import (
	"context"
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
	watcher clientv3.Watcher
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
		watcher clientv3.Watcher
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
	watcher = clientv3.NewWatcher(client)
	G_JobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
		watcher: watcher,
	}

	//启动任务监听
	G_JobMgr.watchJobs()

	return
}

//创建任务执行锁
func (jobMgr *JobMgr) CreateJobLock(jobName string) (jobLock *JobLock) {
	jobLock = &JobLock{
		JobName:    jobName,
		KV:         jobMgr.kv,
		Lease:      jobMgr.lease,
	}
	return
}

// 监听任务变化
func (jobMgr *JobMgr) watchJobs() (err error) {
	var (
		getResp *clientv3.GetResponse
		watchRevision int64
		kv *mvccpb.KeyValue
		job *common.Job
		jobEvent *common.JobEvent
		watchChan clientv3.WatchChan
		watchResp clientv3.WatchResponse
		watchEvent *clientv3.Event
		jobName string
	)
	if getResp, err = jobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}
	for _,kv = range getResp.Kvs {
		if job, err = common.UnpackJob(kv.Value); err == nil {
			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			G_scheduler.PushJobEvent(jobEvent)
		}
	}

	// 监听
	go func() {
		watchRevision = getResp.Header.Revision + 1
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(watchRevision), clientv3.WithPrefix())
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
				case mvccpb.DELETE:
					jobName = string(watchEvent.Kv.Key)
					job = &common.Job{Name: jobName}
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
				}
				G_scheduler.PushJobEvent(jobEvent)
			}
		}
	}()
	return
}






























