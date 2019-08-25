package engine

import (
	"time"
	"windcontrol-go/config/types"
)

type Engine struct {
	Services []types.TickService
}

func (engine Engine) Run() {
	for _, service := range engine.Services {
		execService(service)
	}
}

func execService(service types.TickService)  {
	go func() {
		tick := time.Tick(service.Interval)
		for {
			<- tick
			service.JobHandle()
		}
	}()
}
