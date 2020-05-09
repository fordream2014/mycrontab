package main

import (
	"flag"
	"fmt"
	"mycrontab/worker"
	"runtime"
	"time"
)
var (
	confFile string  //配置文件路径
)
func initArgs() {
	//worker -config ./master.json -xxx 222
	//worker -h
	flag.StringVar(&confFile, "config", "./worker.json", "指定worker.json")
	flag.Parse()
}

//初始化线程数量
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)
	//初始化命令行参数
	initArgs()
	//初始化线程
	initEnv()
	//加载配置
	if err = worker.InitConfig(confFile); err != nil {
		goto END
	}
	// 服务注册
	if err = worker.InitRegister(); err != nil {
		goto END
	}
	//启动执行器
	if err = worker.InitExecutor(); err != nil {
		goto END
	}
	//启动调度器
	if err = worker.InitScheduler(); err != nil {
		goto END
	}
	//任务管理器
	if err = worker.InitJobMgr(); err != nil {
		goto END
	}
	fmt.Println("启动成功")
	//正常退出
	for {
		time.Sleep(1 * time.Second)
	}
END:
	fmt.Println(err)
}
