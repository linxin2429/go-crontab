package master

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
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
		common.HttpInputErrorHandle(w, err)
		return
	}
	oldJob, err := Global_JobMgr.SaveJob(job)
	if err != nil {
		common.HttpInternalErrorHandle(w, err)
		return
	}
	common.Logger.Infof("job save %s", postJob)
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
		return
	}
	name := r.PostForm.Get("name")
	oldJob, err := Global_JobMgr.DeleteJob(name)
	if err != nil {
		common.HttpInputErrorHandle(w, err)
		return
	}
	if oldJob == nil {
		common.HttpInputErrorHandle(w, errors.New(fmt.Sprintf("job name[%s] not found", name)))
		return
	}
	common.Logger.Infof("job delete: %s", name)
	data, err := common.BuildResponse(0, "success", oldJob)
	if err != nil {
		common.HttpInternalErrorHandle(w, err)
		return
	}
	w.Write(data)
}

// GET /job/list
func handleJobList(w http.ResponseWriter, r *http.Request) {
	jobs, err := Global_JobMgr.ListJob()
	common.Logger.Infof("job list")
	if err != nil {
		common.HttpInternalErrorHandle(w, err)
		return
	}
	data, err := common.BuildResponse(0, "success", jobs)
	if err != nil {
		common.HttpInternalErrorHandle(w, err)
		return
	}
	w.Write(data)
}

func handleJobKill(w http.ResponseWriter, r *http.Request) {

}

// InitApiServer 初始化apiserver
func InitApiServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)

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
