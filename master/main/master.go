package main

import (
	"flag"
	"fmt"
	"mycrontab/master"
	"runtime"
	"time"
)
var (
	confFile string  //配置文件路径
)
func initArgs() {
	//master -config ./master.json -xxx 222
	//master -h
	flag.StringVar(&confFile, "config", "./master.json", "指定master.json")
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
	if err = master.InitConfig(confFile); err != nil {
		goto END
	}
	// 初始化服务发现模块
	if err = master.InitWorkerMgr(); err != nil {
		goto END
	}
	//任务管理器
	if err = master.InitJobMgr(); err != nil {
		goto END
	}
	// 启动Api HTTP服务
	if err = master.InitApiServer(); err != nil {
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
