package common

import (
	"context"
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

// 定时任务
type Job struct {
	Name string `json:"name"`	//任务名
	Command string `json:"command"`	//shell命令
	CronExpr string `json:"cronExpr"`	//cron表达式
}
// Http接口应答
type Response struct {
	Errno	int	`json:"errno"`
	Msg	string	`json:"msg"`
	Data	interface{}	`json:"data"`
}
type JobEvent struct {
	EventType int64
	Job *Job
}
// 任务调度计划
type JobSchedulePlan struct {
	Job *Job // 要调度的任务信息
	Expr *cronexpr.Expression //解析好的cronexpr表达式
	NextTime time.Time //下次调度时间
}
//任务执行状态
type JobExecuteInfo struct {
	Job *Job //任务信息
	PlanTime time.Time //理论上的调度时间
	RealTime time.Time //实际的调度时间
	CancelCtx context.Context //任务command的context
	CancelFunc context.CancelFunc //用于取消command执行的cancel函数
}

//任务执行结果
type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo
	Output []byte
	Err error
	StartTime time.Time
	EndTime time.Time
}

//应答方法
func BuildResponse(errno int, msg string, data interface{}) (resp []byte, err error) {
	var (
		response Response
	)
	response.Errno = errno
	response.Msg = msg
	response.Data = data
	//序列化json
	resp, err = json.Marshal(response)
	return
}

func UnpackJob(bytes []byte)(job *Job, err error) {
	job = &Job{}
	if err = json.Unmarshal(bytes, job); err != nil {
		return
	}
	return
}

func BuildJobEvent(eventType int64, job *Job)(*JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job: job,
	}
}

// 构造任务执行计划
func BuildJobSchedulePlan(job *Job) (jobSchedulePlan *JobSchedulePlan, err error) {
	var (
		expr *cronexpr.Expression
	)
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}
	jobSchedulePlan = &JobSchedulePlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}
	return
}
// 构造执行状态信息
func BuildJobExecuteInfo(plan *JobSchedulePlan) (jobExecuteInfo *JobExecuteInfo) {
	var (
		cancelContext context.Context
		cancelFunc context.CancelFunc
	)
	cancelContext, cancelFunc = context.WithCancel(context.TODO())
	jobExecuteInfo = &JobExecuteInfo{
		Job:        plan.Job,
		PlanTime:   plan.NextTime, //计划调度事假
		RealTime:   time.Now(), //真实调度时间
		CancelCtx:  cancelContext,
		CancelFunc: cancelFunc,
	}
	return
}

// 从etcd的key中提取任务名
// /cron/jobs/job10抹掉/cron/jobs/
func ExtractJobName(jobKey string) (string) {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

// 从 /cron/killer/job10提取job10
func ExtractKillerName(killerKey string) (string) {
	return strings.TrimPrefix(killerKey, JOB_KILL_DIR)
}

// 提取worker的IP
func ExtractWorkerIP(regKey string) (string) {
	return strings.TrimPrefix(regKey, JOB_WORKER_DIR)
}











