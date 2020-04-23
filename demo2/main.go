package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	mockCron()
}

type Crontab struct {
	expr *cronexpr.Expression
	next time.Time
}

func mockCron() {
	//遍历crontab，判断是否可以运行
	var (
		cronMap map[string]*Crontab
		now     time.Time
		expr    *cronexpr.Expression
	)
	now = time.Now()
	cronMap = make(map[string]*Crontab)
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronMap["job1"] = &Crontab{
		expr: expr,
		next: expr.Next(now),
	}

	expr = cronexpr.MustParse("*/2 * * * * * *")
	cronMap["job2"] = &Crontab{
		expr: expr,
		next: expr.Next(now),
	}

	go func() {
		var (
			now     time.Time
			jobName string
			jobCron *Crontab
		)
		for {
			now = time.Now()
			for jobName, jobCron = range cronMap {
				if jobCron.next.Before(now) || jobCron.next.Equal(now) {
					go func(jobName string) {
						fmt.Println(jobName, "执行")
					}(jobName)

					//设置下次执行时间
					jobCron.next = jobCron.expr.Next(now)
					fmt.Println(jobName, "下次执行时间为:", jobCron.next)
				}
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()

	time.Sleep(100 * time.Second)
}

func testCron() {
	var (
		expr    *cronexpr.Expression
		err     error
		current time.Time
		next    time.Time
	)
	//cronexpr 秒粒度
	if expr, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		fmt.Println(err)
		return
	}

	current = time.Now()
	next = expr.Next(current)

	time.AfterFunc(next.Sub(current), func() {
		fmt.Println("被调度了：", next)
	})
	time.Sleep(10 * time.Second)
}
