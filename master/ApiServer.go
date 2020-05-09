package master

import (
	"encoding/json"
	"fmt"
	"mycrontab/common"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

var (
	G_apiServer *ApiServer
)

//保存任务接口
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	var(
		err error
		postJob string
		job common.Job
		oldJob *common.Job
		bytes []byte
	)
	//解析post表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	//获取表单中的job字段
	postJob = req.PostForm.Get("job")
	//postJob = `{"name":"job2", "command":"echo hello job2","cronExpr":"* * * * *"}`
	fmt.Println(postJob)
	//反序列化为json
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	//保存到ETCD中
	if oldJob,err = G_JobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	//返回正常应答
	if bytes,err = common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}
//删除任务
func handleJobDelete(resp http.ResponseWriter, req *http.Request) {
	var (
		err error
		jobName string
		oldJob *common.Job
		bytes []byte
	)
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	jobName = req.PostForm.Get("name")
	if oldJob, err = G_JobMgr.DeleteJob(jobName); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}
//删除任务
func handleJobList(resp http.ResponseWriter, req *http.Request) {
	var (
		err error
		jobs []*common.Job
		bytes []byte
	)
	if jobs, err = G_JobMgr.GetJobList(); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(0, "success", jobs); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

//杀死任务
func handleJobKill(resp http.ResponseWriter, req *http.Request) {
	var (
		err error
		bytes []byte
		jobName string
	)
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	jobName = req.PostForm.Get("name")
	if jobName == "" {
		err = fmt.Errorf("param invalid!")
		goto ERR
	}
	if err = G_JobMgr.JobKill(jobName); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(0, "success", nil); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}
//获取健康worker节点列表
func handleWorkerList(resp http.ResponseWriter, req *http.Request) {
	var (
		workerArr []string
		err error
		bytes []byte
	)

	if workerArr, err = G_workerMgr.ListWorkers(); err != nil {
		goto ERR
	}

	// 正常应答
	if bytes, err = common.BuildResponse(0, "success", workerArr); err == nil {
		resp.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func InitApiServer() (err error){
	var (
		mux *http.ServeMux
		listener net.Listener
		httpServer *http.Server
		staticDir http.Dir
		staticHandler http.Handler
	)
	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)
	mux.HandleFunc("/worker/list", handleWorkerList)

	staticDir = http.Dir(G_config.WebRoot)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))

	//启动TCP监听
	if listener, err = net.Listen("tcp", ":" + strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}
	//创建一个HTTP服务
	httpServer = &http.Server{
		Handler:           mux,
		ReadTimeout:       time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout:      time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
	}
	//赋值单例
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}
	go httpServer.Serve(listener)
	return
}
