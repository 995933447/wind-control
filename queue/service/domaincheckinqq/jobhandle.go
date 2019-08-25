package domaincheckinqq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"util/sms"
	"util/timer"
	"util/url"
	"windcontrol-go/config"
	"windcontrol-go/config/types"
	"windcontrol-go/logger"
	"windcontrol-go/model"
	domainCheckInQqConfig "windcontrol-go/queue/service/domaincheckinqq/config"
)

var domainContent = make(map[string][]byte)

type tokenResponse struct {
	Data tokenData
}

type tokenData struct {
	Timestamp int
	Token string
}

type checkDomainResponse struct {
	Data domainData
}

type domainData struct {
	Domain string
	Intercept int
}

const (
	stop = "停用"
)

func JobHandle(task types.Task) error {
	domain := task.(model.Domains)
	interval := time.Duration(domain.QqCheckInterval)
	time.Sleep(time.Minute * interval)
	checkDomain(domain)
	tasks.Store(domain, freeStatus)
	return nil
}

func checkDomain(domain model.Domains) {
	logger.DefaultLogger.Debug("Checking url:" + domain.Link + " in qq", nil)

	now := timer.NowUnix(nil)

	tokenResp, err := getTokenResp(&domain)
	if err != nil {
		logger.DefaultLogger.Error(err, nil)
		return
	}

	checkDomainResp, err := getCheckDomainResp(&domain, tokenResp)
	if err != nil {
		logger.DefaultLogger.Error(err, nil)
		return
	}

	domain.QqStatus = model.DomainQqValidStatus
	if checkDomainResp.Data.Intercept == domainCheckInQqConfig.DomainInValidStatus {
		domain.QqStatus = model.DomainQqInvalidStatus
		domain.QqStopTime = now
		if err := notifyWarn(domain.Link, stop); err != nil {
			logger.DefaultLogger.Error(fmt.Sprintf("sending mail err: %v", err), nil)
		}
	} else if checkDomainResp.Data.Intercept == domainCheckInQqConfig.DomainCheckFailedStatus {
		logger.DefaultLogger.Error(checkDomainResp, nil)
	}
	fmt.Printf("%d %+v\n", checkDomainResp.Data.Intercept, domain)

	domain.LastCheckTime = now
	row, err := mysqlClient.Id(domain.Id).Cols(model.DomainQqStatusField, model.DomainQqStopTimeField, model.DomainLastCheckTimeFiled).Update(&domain)
	if err != nil {
		logger.DefaultLogger.Error(fmt.Sprintf("updating database error \n", err), nil)
	}
	logger.DefaultLogger.Info(fmt.Sprintf("%d updated", row), nil)
}

func getTokenResp(domain *model.Domains) (tokenResponse, error) {
	domain.QqStatus = model.DomainQqValidStatus

	resp, err := http.Get(domainCheckInQqConfig.GetTokenApi)
	if err != nil {
		return tokenResponse{}, err
	}

	tokenResult, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var tokenResp tokenResponse
	err = json.Unmarshal(tokenResult, &tokenResp)
	return tokenResp, err
}

func getCheckDomainResp(domain *model.Domains, tokenResp tokenResponse) (checkDomainResponse, error) {
	host, _ := url.GetHost(domain.Link)
	resp, err := http.Post(domainCheckInQqConfig.CheckDomainApi,
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("timestamp=%d&token=%s&domain=%s", tokenResp.Data.Timestamp, tokenResp.Data.Token, host)))

	if err != nil {
		return checkDomainResponse{}, err
	}

	checkResult, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return checkDomainResponse{}, err
	}


	var checkDomainResp checkDomainResponse
	err = json.Unmarshal(checkResult, &checkDomainResp)
	return checkDomainResp, err
}

func notifyWarn(url string, reason string) error {
	msg := fmt.Sprintf("%s 检测到 %s 被 %s", timer.NowDate("Y-m-d H:i:s"), url, reason)
	subject := fmt.Sprintf("%s异常", url)
	mail := sms.NewEmail(config.ToMail, subject, msg)
	return  sms.SendEmail(mail, config.MailUsername, config.MailPassword, config.MailHost, fmt.Sprintf("%s:%d", config.MailHost, config.MailPort))
}