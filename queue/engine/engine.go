package engine

import (
	"windcontrol-go/config/types"
	"windcontrol-go/queue/schedule"
)

type Engine struct {
	Schedule *schedule.Schedule
}

func (engine *Engine) Run() {
	finishChan := make(types.FinishChan)

	for index, service := range engine.Schedule.Services {
		taskChan := make(chan types.Task)
		engine.Schedule.Services[index].Connection.Popper = taskChan
		failure := service.FailedJobHandle(finishChan)
		for i := 0; i < service.WorkerNum; i++ {
			createWorker(taskChan, failure, service.JobHandle)
		}
	}

	engine.Schedule.Run()

	for {
		<- finishChan
	}
}

func createWorker(taskChan chan types.Task, failureChan chan error, jobHandle types.JobHandle)  {
	go func() {
		for {
			task := <- taskChan
			err := jobHandle(task)
			failureChan <- err
		}
	}()
}