package queue

import (
	"windcontrol-go/config/types"
	"windcontrol-go/queue/engine"
	"windcontrol-go/queue/schedule"
	"windcontrol-go/queue/service/createserverjs"
	domainCheck "windcontrol-go/queue/service/domaincheck"
	domainCheckInQq "windcontrol-go/queue/service/domaincheckinqq"
	domainCheckInWechat "windcontrol-go/queue/service/domaincheckinwechat"
)

func Run()  {
	s := schedule.Schedule{
		registerQueueService(),
	}
	e := engine.Engine{
		Schedule: &s,
	}
	(&e).Run()
}

func registerQueueService() []types.QueueService {
	var services []types.QueueService

	services = append(services,
		types.QueueService{
			types.MakeConnection(domainCheck.Connection),
			100,
			domainCheck.JobHandle,
			domainCheck.FailedJobHandle,
		},
		types.QueueService{
			types.MakeConnection(createserverjs.Connection),
			100,
			createserverjs.JobHandle,
			createserverjs.FailedJobHandle,
		},
		types.QueueService{
			types.MakeConnection(domainCheckInQq.Connection),
			100,
			domainCheckInQq.JobHandle,
			domainCheckInQq.FailedJobHandle,
		},
		types.QueueService{
			types.MakeConnection(domainCheckInWechat.Connection),
			100,
			domainCheckInWechat.JobHandle,
			domainCheckInWechat.FailedJobHandle,
		},
	)

	return services
}
