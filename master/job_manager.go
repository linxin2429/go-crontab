package master

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
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

// SaveJob 保存任务
func (jobMgr *JobMgr) SaveJob(job *common.Job) (*common.Job, error) {
	jobKey := common.ETCD_JOB_SAVE_DIR + job.Name
	jobValue, err := json.Marshal(job)
	if err != nil {
		return nil, errors.Wrap(err, "json marshal error")
	}
	putResponse, err := jobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV())
	if err != nil {
		return nil, errors.Wrap(err, "etcd put error")
	}
	if putResponse.PrevKv != nil {
		oldKv := &common.Job{}
		err := json.Unmarshal(putResponse.PrevKv.Value, oldKv)
		return oldKv, errors.Wrap(err, "json unmarshal err")
	}
	return nil, nil
}

func (jobMgr *JobMgr) DeleteJob(name string) (*common.Job, error) {
	jobKey := common.ETCD_JOB_SAVE_DIR + name
	deleteResp, err := jobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV())
	if err != nil {
		err = errors.Wrap(err, "etcd delete error")
		common.Logger.Errorln(err)
		return nil, err
	}
	if len(deleteResp.PrevKvs) != 0 {
		oldKv := &common.Job{}
		err := json.Unmarshal(deleteResp.PrevKvs[0].Value, oldKv)
		return oldKv, errors.Wrap(err, "json unmarshal error")
	}
	return nil, nil
}

func (jobMgr *JobMgr) ListJob() ([]*common.Job, error) {
	getResponse, err := jobMgr.kv.Get(context.TODO(), common.ETCD_JOB_SAVE_DIR, clientv3.WithPrefix())
	if err != nil {
		return nil, errors.Wrap(err, "etcd get error")
	}
	var jobs = make([]*common.Job, 0)
	for _, r := range getResponse.Kvs {
		job := &common.Job{}
		err := json.Unmarshal(r.Value, job)
		if err != nil {
			continue
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}
