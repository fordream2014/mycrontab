package worker

import (
	"fmt"
	"mycrontab/common"
	"time"
)

type Scheduler struct {
	jobEventChan chan *common.JobEvent
	jobPlanTable map[string]*common.JobSchedulePlan
	jobExecutingTable map[string]*common.JobExecuteInfo
	jobResultChan chan *common.JobExecuteResult
}

var (
	G_scheduler *Scheduler
)

func (scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchedulePlan) {
	var (
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting bool
	)
	//如果任务正在执行，跳过本次调度
	if jobExecuteInfo, jobExecuting = scheduler.jobExecutingTable[jobPlan.Job.Name]; jobExecuting {
		fmt.Println("尚未执行结束，", jobExecuteInfo.Job.Name)
		return
	}
	jobExecuteInfo = common.BuildJobExecuteInfo(jobPlan)

	//保存执行状态
	scheduler.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo

	//任务执行
	G_executor.ExecuteJob(jobExecuteInfo)
}

func (scheduler *Scheduler) TrySchedule() (scheduleAfter time.Duration){
	var (
		jobSchedulePlan *common.JobSchedulePlan
		now time.Time
		nearTime *time.Time
	)
	if len(scheduler.jobPlanTable) == 0 {
		scheduleAfter = 1 * time.Second;
		return
	}
	now = time.Now()
	for _, jobSchedulePlan = range scheduler.jobPlanTable {
		if jobSchedulePlan.NextTime.Before(time.Now()) || jobSchedulePlan.NextTime.Equal(time.Now()) {
			// 更新下次执行时间
			jobSchedulePlan.NextTime = jobSchedulePlan.Expr.Next(now)
			// 执行任务
			scheduler.TryStartJob(jobSchedulePlan)
		}

		if nearTime == nil || jobSchedulePlan.NextTime.Before(*nearTime) {
			nearTime = &jobSchedulePlan.NextTime
		}
	}
	//下次调度时间
	scheduleAfter = (*nearTime).Sub(now)
	return
}

//处理任务事件
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		jobSchedulePlan *common.JobSchedulePlan
		jobExecuteInfo *common.JobExecuteInfo
		err error
		jobExisted bool
		jobExecuting bool
	)
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE: //保存任务
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		scheduler.jobPlanTable[jobSchedulePlan.Job.Name] = jobSchedulePlan
	case common.JOB_EVENT_DELETE:
		if jobSchedulePlan, jobExisted = scheduler.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	case common.JOB_EVENT_KILL:
		if jobExecuteInfo, jobExecuting = scheduler.jobExecutingTable[jobEvent.Job.Name]; jobExecuting {
			jobExecuteInfo.CancelFunc()
		}
	}
}

// 调度协程
func (scheduler *Scheduler) scheduleLoop() {
	var (
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
		jobEvent *common.JobEvent
		jobResult *common.JobExecuteResult
	)
	//初始化
	scheduleAfter = scheduler.TrySchedule()
	fmt.Println("调度延迟(秒)：", scheduleAfter.Seconds())
	//调度延迟定时器
	scheduleTimer = time.NewTimer(scheduleAfter)

	//定时任务
	for {
		select {
		case jobEvent = <-scheduler.jobEventChan :
			scheduler.handleJobEvent(jobEvent)
		case <-scheduleTimer.C:
		case jobResult = <-scheduler.jobResultChan:
			scheduler.handleJobResult(jobResult)
		}

		scheduleAfter = scheduler.TrySchedule()
		fmt.Println("调度延迟(秒)：", scheduleAfter.Seconds())
		scheduleTimer.Reset(scheduleAfter)
	}
}

func (scheduler *Scheduler) handleJobResult(result *common.JobExecuteResult) {
	delete(scheduler.jobExecutingTable, result.ExecuteInfo.Job.Name)
	fmt.Println("任务执行完成 ", result.ExecuteInfo.Job.Name, string(result.Output), result.Err)
}

func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	fmt.Println("push job event", jobEvent.Job.Name)
	scheduler.jobEventChan <- jobEvent
}

func (scheduler *Scheduler) PushJobResult(jobResult *common.JobExecuteResult) {
	scheduler.jobResultChan <- jobResult
}

func InitScheduler()(err error) {
	G_scheduler = &Scheduler{
		jobEventChan:      make(chan *common.JobEvent, 1000),
		jobPlanTable:      make(map[string]*common.JobSchedulePlan),
		jobExecutingTable: make(map[string]*common.JobExecuteInfo),
		jobResultChan:     make(chan *common.JobExecuteResult, 1000),
	}
	go G_scheduler.scheduleLoop()
	return
}