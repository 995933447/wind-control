package createserverjs

import (
	"windcontrol-go/config/types"
)

func FailedJobHandle(finishChan types.FinishChan) chan error {
	errChan := make(chan error)
	go func() {
		for {
			err := <- errChan
			if err != nil {
				panic(err)
			}
			finishChan <- true
		}
	}()
	return errChan
}
