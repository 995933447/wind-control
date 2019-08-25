package domaincheck

import (
	"fmt"
	"net/smtp"
	"strings"
	"testing"
)

const (
	HOST        = "smtp.163.com"
	SERVER_ADDR = "smtp.163.com:25"
	USER        = "CYLzy12@163.com" //发送邮件的邮箱
	PASSWORD    = "aa123456"         //发送邮件邮箱的密码
)

func TestJobHandle(t *testing.T) {
	mycontent := " my dear令"

	email := NewEmail("995933447@qq.com;", "test golang email", mycontent)
	err := SendEmail(email)

	fmt.Println(err)
}

type Email struct {
	to      string "to"
	subject string "subject"
	msg     string "msg"
}

func NewEmail(to, subject, msg string) *Email {
	return &Email{to: to, subject: subject, msg: msg}
}

func SendEmail(email *Email) error {
	auth := smtp.PlainAuth("", USER, PASSWORD, HOST)
	sendTo := strings.Split(email.to, ";")
	done := make(chan error, 1024)

	go func() {
		defer close(done)
		for _, v := range sendTo {

			str := strings.Replace("From: "+USER+"~To: "+v+"~Subject: "+email.subject+"~~", "~", "\r\n", -1) + email.msg

			err := smtp.SendMail(
				SERVER_ADDR,
				auth,
				USER,
				[]string{v},
				[]byte(str),
			)
			done <- err
		}
	}()

	for i := 0; i < len(sendTo); i++ {
		<-done
	}

	return nil
}