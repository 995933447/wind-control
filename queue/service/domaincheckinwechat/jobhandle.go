package domaincheckinwechat

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
	domainCheckInWechatConfig "windcontrol-go/queue/service/domaincheckinwechat/config"
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
	interval := time.Duration(domain.WechatCheckInterval)
	time.Sleep(time.Minute * interval)
	checkDomain(domain)
	tasks.Store(domain, freeStatus)
	return nil
}

func checkDomain(domain model.Domains) {
	logger.DefaultLogger.Debug("Checking url:" + domain.Link + " in wechat", nil)

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

	if checkDomainResp.Data.Intercept == domainCheckInWechatConfig.DomainInValidStatus {
		domain.WechatStatus = model.DomainWechatInvalidStatus
		domain.WechatStopTime = now
		if err := notifyWarn(domain.Link, stop); err != nil {
			logger.DefaultLogger.Error(fmt.Sprintf("sending mail err: %v", err), nil)
		}
	} else if checkDomainResp.Data.Intercept == domainCheckInWechatConfig.DomainCheckFailedStatus {
		logger.DefaultLogger.Error(checkDomainResp, nil)
	}

	domain.LastCheckTime = now
	_, err = mysqlClient.Id(domain.Id).Cols(model.DomainWechatStatusField, model.DomainWechatStopTimeField, model.DomainLastCheckTimeFiled).Update(&domain)
	if err != nil {
		logger.DefaultLogger.Error(fmt.Sprintf("updating database error \n", err), nil)
	}
}

func getTokenResp(domain *model.Domains) (tokenResponse, error) {
	domain.WechatStatus = model.DomainWechatValidStatus

	resp, err := http.Get(domainCheckInWechatConfig.GetTokenApi)
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
	resp, err := http.Post(domainCheckInWechatConfig.CheckDomainApi,
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