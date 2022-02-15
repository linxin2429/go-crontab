package common

import (
	"encoding/json"
)

// Job 保存至ETCD的任务
type Job struct {
	// 任务名
	Name string `json:"name"`
	// 任务指令
	Command string `json:"command"`
	// cron表达式
	CronExpr string `json:"cron_expr"`
}

// HTTP response
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

type ErrorResponse struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
}

func BuildResponse(errno int, msg string, data interface{}) ([]byte, error) {
	response := &Response{
		Errno: errno,
		Msg:   msg,
		Data:  data,
	}
	return json.Marshal(response)
}
func BuildErrorResponse(errno int, err error) []byte {
	response := &ErrorResponse{
		Errno: errno,
		Msg:   err.Error(),
	}
	data, _ := json.Marshal(response)
	return data
}
