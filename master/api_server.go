package master

import (
	"encoding/json"
	"fmt"
	"go-crontab/common"
	"net"
	"net/http"
	"time"
)

// Global_ApiServer 全局Api服务器
var Global_ApiServer *ApiServer

// ApiServer 任务的HTTP接口
type ApiServer struct {
	httpServer *http.Server
}

// 保存任务
// POST /job/save
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		common.HttpInternalErrorHandle(w, err)
		return
	}
	postJob := r.PostForm.Get("job")
	job := &common.Job{}
	err = json.Unmarshal([]byte(postJob), job)
	if err != nil {
		common.HttpInternalErrorHandle(w, err)
		return
	}
	oldJob, err := Global_JobMgr.SaveJob(job)
	if err != nil {
		common.HttpInternalErrorHandle(w, err)
		return
	}
	data, err := common.BuildResponse(0, "success", oldJob)
	if err != nil {
		common.HttpInternalErrorHandle(w, err)
		return
	}
	w.Write(data)
}

// POST /job/delete
func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		common.HttpInternalErrorHandle(w, err)
	}
	name := r.PostForm.Get("name")
	oldJob, err := Global_JobMgr.DeleteJob(name)
	if err != nil {
		common.HttpInternalErrorHandle(w, err)
	}
	data, err := common.BuildResponse(0, "success", oldJob)
	if err != nil {
		common.HttpInternalErrorHandle(w, err)
		return
	}
	w.Write(data)
}

// InitApiServer 初始化apiserver
func InitApiServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", Global_Config.ApiPort))
	if err != nil {
		return err
	}
	httpServer := &http.Server{
		ReadTimeout:  time.Duration(Global_Config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(Global_Config.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	Global_ApiServer = &ApiServer{
		httpServer: httpServer,
	}

	go httpServer.Serve(listener)

	return nil
}
