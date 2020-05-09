package worker

import (
	"mycrontab/common"
	"os/exec"
	"time"
)

type Executor struct {}
var (
	G_executor *Executor
)

func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		var(
			result *common.JobExecuteResult
			cmd *exec.Cmd
			output []byte
			err error
			jobLock *JobLock
		)
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}
		//初始化分布式锁
		jobLock = G_JobMgr.CreateJobLock(info.Job.Name)
		err = jobLock.TryLock()
		defer jobLock.Unlock()

		if err != nil {
			//上锁失败
			result.Err = err
			result.EndTime = time.Now()
		} else {
			result.StartTime = time.Now()
			cmd = exec.CommandContext(info.CancelCtx, "/bin/bash", "-c", info.Job.Command)
			output, err = cmd.CombinedOutput()

			result.EndTime = time.Now()
			result.Err = err
			result.Output = output
		}
		//任务执行完成，把结果返回给schedulor，schedulor从executingTable中删除掉执行记录
		G_scheduler.PushJobResult(result)
	}()
}

func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}
