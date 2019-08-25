package domaincheckinwechat

import (
	"windcontrol-go/config/types"
	"windcontrol-go/logger"
)

func FailedJobHandle(finishChan types.FinishChan) chan error {
	errChan := make(chan error)
	go func() {
		for {
			err := <- errChan
			if err != nil {
				logger.DefaultLogger.Error(err, nil)
			}
			finishChan <- true
		}
	}()
	return errChan
}
