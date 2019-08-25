package domaincheckinwechat

import (
	"sync"
	"util/timer"
	"windcontrol-go/config/types"
	"windcontrol-go/logger"
	"windcontrol-go/model"
)

const (
	freeStatus  = 0
	busyStatus = 1
)

var lastFreshTasksAt = 0

var freeTasks = make([]model.Domains, 0)

var tasks sync.Map

func Connection() (task types.Task, has bool) {
	if len(freeTasks) <= 0 {
		err := freshTasks()
		if err != nil {
			logger.DefaultLogger.Error(err, nil)
		}
		getFreeTasks()
		if len(freeTasks) <= 0 {
			return nil, false
		}
	}

	task = freeTasks[0]
	freeTasks = freeTasks[1:]
	return task, true
}

func getFreeTasks() []model.Domains {
	tasks.Range(func(task, status interface{}) bool {
		if status == freeStatus {
			freeTasks = append(freeTasks, task.(model.Domains))
			tasks.Store(task, busyStatus)
		}
		return true
	})

	return freeTasks
}

func freshTasks() error {
	now := timer.NowUnix(nil)
	if (now - lastFreshTasksAt) < 60 {
		return nil
	}

	var domains []model.Domains
	err := mysqlClient.
		Where(model.DomainWechatStatusField + " = ?", model.DomainWechatValidStatus).
		Where(model.DomainIsEnableWechatCheckField + " = ?", model.DomainIsEnableWechatCheck).
		Cols(model.DomainWechatCheckIntervalField, model.DomainLinkField, model.DomainIdField).
		Find(&domains)
	if err != nil {
		return err
	}

	lastFreshTasksAt = now

	for _, domain := range domains {
		if _, ok := tasks.Load(domain); !ok {
			tasks.Store(domain, freeStatus)
		}
	}

	tasks.Range(func(task, status interface{}) bool {
		has := false
		for _, domain := range domains {
			if domain == task.(model.Domains) {
				has = true
			}
		}
		if !has {
			tasks.Delete(task)
		}
		return true
	})

	return nil
}

