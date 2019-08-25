package model

import "time"

const (
	DomainIdField = "id"
	DomainStatusField = "status"
	DomainQqStatusField = "qq_status"
	DomainWechatStatusField = "wechat_status"
	DomainStopTimeField = "stop_time"
	DomainQqStopTimeField = "qq_stop_time"
	DomainWechatStopTimeField = "wechat_stop_time"
	DomainLinkField = "link"
	DomainCheckIntervalField = "check_interval"
	DomainFromField = "from"
	DomainLastCheckTimeFiled = "last_check_time"
	DomainIsEnableQqCheckField = "is_enable_qq_check"
	DomainIsEnableWechatCheckField = "is_enable_wechat_check"
	DomainQqCheckIntervalField = "qq_check_interval"
	DomainWechatCheckIntervalField = "wechat_check_interval"

	DomainValidStatus = 1
	DomainInvalidStatus = 0
	DomainQqValidStatus = 1
	DomainQqInvalidStatus = 0
	DomainWechatValidStatus = 1
	DomainWechatInvalidStatus = 0
	DomainDomainType = 0
	DomainIpType = 1
	DomainFromApiType = 0
	DomainFromStaticType = 1
	DomainFromLoginType = 2
	DomainFromFrontType = 3
	DomainIsEnableQqCheck = 1
	DomainIsNotEnableQqCheck = 0
	DomainIsEnableWechatCheck = 1
	DomainIsNotEnableWechatCheck = 0
)

type Domains struct {
	Id int64
	Link string
	Type int
	Status int
	QqStatus int
	WechatStatus int
	IsEnableQqCheck int
	IsEnableWechatCheck int
	CheckInterval int
	QqCheckInterval int
	WechatCheckInterval int
	From int
	StopTime int
	QqStopTime int
	WechatStopTime int
	LastCheckTime int
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

func (Domains) TableName() string {
	return "domains"
}