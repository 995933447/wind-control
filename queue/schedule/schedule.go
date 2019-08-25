package schedule

import (
	"windcontrol-go/config/types"
)

type Schedule struct {
	Services []types.QueueService
}

func (schedule *Schedule) Run() {
	go func() {
		for _, service := range schedule.Services {
			popTask(service)
		}
	}()
}

func popTask(service types.QueueService) {
	go func() {
		for {
			task, has := service.Connection.PopTask()
			if has {
				service.Connection.Popper <- task
			}
		}
	}()
}