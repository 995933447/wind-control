package domaincheck

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

var tasks sync.Map

var freeTasks = make([]model.Domains, 0)

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
			tasks.Store(task, busyStatus)
			freeTasks = append(freeTasks, task.(model.Domains))
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
	err := mysqlClient.Where(model.DomainStatusField + " = ?", model.DomainValidStatus).Cols(model.DomainLinkField, model.DomainCheckIntervalField).Find(&domains)
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

