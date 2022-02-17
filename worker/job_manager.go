package worker

import (
	"context"
	"go-crontab/common"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// JobMgr 任务管理类
type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

// Global_JobMgr 全局任务管理类
var Global_JobMgr *JobMgr

// InitJobMgr 初始化
func InitJobMgr() error {
	config := clientv3.Config{
		Endpoints:   Global_Config.EtcdEndpoints,
		DialTimeout: time.Duration(Global_Config.EtcdDialTimeout) * time.Millisecond,
		Username:    Global_Config.EtcdUsername,
		Password:    Global_Config.EtcdPwd,
	}
	client, err := clientv3.New(config)
	if err != nil {
		return err
	}
	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)

	Global_JobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return nil
}

func (jobMgr *JobMgr) watchJobs() error {
	getResponse, err := jobMgr.kv.Get(context.TODO(), common.ETCD_JOB_SAVE_DIR, clientv3.WithPrefix())
	if err != nil {
		return err
	}

}
