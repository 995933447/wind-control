package config

import "fmt"

type QueueService struct {
	Connection Connection
	WorkerNum int
	JobHandle JobHandle
	FailedJobHandle func(FinishChan) chan error
}

type FinishChan chan bool

type Connection struct {
	Popper chan Task
	PopTask func() (task Task, has bool)
}

type JobHandle func(Task) error

type Task interface {}

func MakeNilQueueService() QueueService {
	var nilQueueService QueueService
	nilQueueService.Connection = MakeNilConnection()
	nilQueueService.WorkerNum = 1
	nilQueueService.FailedJobHandle = NilFailedJobHandle
	nilQueueService.JobHandle = NilJobHandle
	return nilQueueService
}

func NilJobHandle(Task) error {
	fmt.Println("Starting job.")
	return nil
}

func MakeNilConnection() Connection {
	var nilConnection Connection
	nilConnection.PopTask = func() (task Task, has bool) {
		fmt.Println("Pop nil connection")
		return nil, true
	}
	return nilConnection
}

func NilFailedJobHandle(finishChan FinishChan) chan error {
	failureChan := make(chan error)
	go func() {
		for {
			<- failureChan
			fmt.Println("Failed job start.")
			finishChan <- true
		}
	}()
	return failureChan
}