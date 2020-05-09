package worker

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"mycrontab/common"
)

type JobLock struct {
	JobName string
	KV clientv3.KV
	Lease clientv3.Lease
	CancelFunc context.CancelFunc
	LeaseId clientv3.LeaseID
	IsLocked bool //是否上锁成功
}

//获取分布式锁
func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (lock *JobLock){
	lock = &JobLock{
		JobName:       jobName,
		KV:            kv,
		Lease:         lease,
	}
	return
}

//txn
func (jobLock *JobLock) TryLock() (err error) {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseKeepAliveRespChan <-chan *clientv3.LeaseKeepAliveResponse
		cancelContext context.Context
		cancelFunc context.CancelFunc
		txn clientv3.Txn
		lockKey string
		txnResp *clientv3.TxnResponse
	)
	//创建1秒的租约
	if leaseGrantResp, err = jobLock.Lease.Grant(context.TODO(), 1); err != nil {
		return
	}
	//创建cancelContext
	cancelContext, cancelFunc = context.WithCancel(context.TODO())

	//自动续租
	if leaseKeepAliveRespChan, err = jobLock.Lease.KeepAlive(cancelContext, leaseGrantResp.ID); err != nil {
		return
	}
	//处理续约应答协程
	go func() {
		var leaseKeepAliveResp *clientv3.LeaseKeepAliveResponse
		for {
			select {
			case leaseKeepAliveResp = <- leaseKeepAliveRespChan:
				if leaseKeepAliveResp == nil {
					goto END
				}
			}
		}
		END:
	}()
	//创建事务txn
	txn = jobLock.KV.Txn(context.TODO())

	//锁路径
	lockKey = common.JOB_LOCK_DIR + jobLock.JobName

	//事务抢锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseGrantResp.ID))).
		Else(clientv3.OpGet(lockKey))

	//提交事务
	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
	}

	//成功返回，失败释放租约
	if !txnResp.Succeeded {
		//锁被占用
		err = common.ERR_LOCK_ALREADY_REQUIRED
		goto FAIL
	}

	jobLock.CancelFunc = cancelFunc
	jobLock.LeaseId = leaseGrantResp.ID
	jobLock.IsLocked = true
	return
FAIL:
	cancelFunc()
	jobLock.Lease.Revoke(context.TODO(), leaseGrantResp.ID)
	return
}

func (jobLock *JobLock) Unlock() {
	if jobLock.IsLocked {
		jobLock.CancelFunc() //取消自动续租协程
		jobLock.Lease.Revoke(context.TODO(), jobLock.LeaseId) //释放租约
	}
}



















