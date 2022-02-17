package main

import (
	"flag"
	"go-crontab/common"
	"go-crontab/worker"
	"runtime"
	"time"
)

// 初始化线程数
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// config配置文件名
var configFile string

// 初始化参数
func initArgs() {
	// master -config ./master.json
	flag.StringVar(&configFile, "config", "./worker.json", "指定worker.json")
	flag.Parse()
}

func main() {
	initArgs()
	initEnv()
	err := worker.InitConfig(configFile)
	common.InitLogger(worker.Global_Config.LogFilename)
	if err != nil {
		common.Logger.Fatalf("error: %s", err)
	}
	err = worker.InitJobMgr()
	if err != nil {
		common.Logger.Fatalf("error: %s", err)
	}
	time.Sleep(10 * time.Minute)
}
