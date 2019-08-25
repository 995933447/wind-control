package createserverjs

import (
	"fmt"
	"time"
	"windcontrol-go/config/types"
	"windcontrol-go/logger"
	"windcontrol-go/model"
	"windcontrol-go/queue/service/createserverjs/config"
)

var tick = time.Tick(config.Interval)

var domainServerTypeMap = map[string]map[int]string {
	"server.js": {
		0: "serverApi",
		1: "serverStaticApi",
		2: "LoginApi",
	},
	"cnzz.js": {
		3: "webUrl",
	},
	"qq.js": {
		4: "qqUrl",
	},
}

func Connection() (types.Task, bool) {
	<- tick

	var result = make(map[string]map[string]string)

	for filename, servers := range domainServerTypeMap {
		addresses := make(map[string]string)
		for index, server := range servers {
				var domain model.Domains

				var has bool
				var err error
				if filename == "qq.js" {
					has, err = mysqlClient.
						Where( model.Domains{}.TableName() + "." + model.DomainFromField + " = ?", index).
						Where(model.DomainQqStatusField + " = ?", model.DomainQqValidStatus).
						Where(model.DomainIsEnableQqCheckField + " = ?", model.DomainIsEnableQqCheck).
						Where(model.DomainStatusField + " = ?", model.DomainValidStatus).
						Get(&domain)
				} else {
					has, err = mysqlClient.
						Where(model.Domains{}.TableName() + "." + model.DomainFromField + " = ?", index).
						Where(model.DomainStatusField + " = ?", model.DomainValidStatus).
						Get(&domain)
				}

				if err != nil {
					logger.DefaultLogger.Error(fmt.Sprintf("select data err: %s", err), nil)
				} else {
					if has {
						addresses[server] = domain.Link
					}
				}

				if _, ok := addresses[server]; !ok {
					addresses[server] = ""
				}
		}
		result[filename] = addresses
	}

	return result, true
}