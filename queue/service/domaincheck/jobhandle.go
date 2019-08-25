package domaincheck

import (
	"fmt"
	"net/http"
	"time"
	"util/sms"
	"util/timer"
	"windcontrol-go/config"
	"windcontrol-go/config/types"
	"windcontrol-go/logger"
	"windcontrol-go/model"
)

var domainContent = make(map[string][]byte)

const (
	hack = "劫持"
	stop = "停用"
)

func JobHandle(task types.Task) error {
	domain := task.(model.Domains)
	interval := time.Duration(domain.CheckInterval)
	time.Sleep(time.Minute * interval)
	checkDomain(domain)
	tasks.Store(domain, freeStatus)
	return nil
}

func checkDomain(domain model.Domains) {
	logger.DefaultLogger.Debug("Checking url:" + domain.Link, nil)

	now := timer.NowUnix(nil)
	_, err := http.Get(domain.Link)

	domain.Status = model.DomainValidStatus

	if err != nil {
		_, err = http.Get(domain.Link)
		if err != nil {
			domain.Status = model.DomainInvalidStatus
			domain.StopTime = now
			logger.DefaultLogger.Error(fmt.Sprintf("Url %s is invalid. error: %v", domain.Link, err), nil)

			if err := notifyWarn(domain.Link, stop); err != nil {
				logger.DefaultLogger.Error(fmt.Sprintf("sending mail err: %v", err), nil)
			}
		}
	}

	domain.LastCheckTime = now
	_, err = mysqlClient.Id(domain.Id).Cols(model.DomainStatusField, model.DomainStopTimeField, model.DomainLastCheckTimeFiled).Update(&domain)
	if err != nil {
		logger.DefaultLogger.Error(fmt.Sprintf("updating database error %s\n", err), nil)
	}
}

func notifyWarn(url string, reason string) error {
	msg := fmt.Sprintf("%s 检测到 %s 被 %s", timer.NowDate("Y-m-d H:i:s"), url, reason)
	subject := fmt.Sprintf("%s异常", url)
	mail := sms.NewEmail(config.ToMail, subject, msg)
	return  sms.SendEmail(mail, config.MailUsername, config.MailPassword, config.MailHost, fmt.Sprintf("%s:%d", config.MailHost, config.MailPort))
}