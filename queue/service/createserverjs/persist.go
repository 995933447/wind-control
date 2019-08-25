package createserverjs

import (
	"github.com/go-xorm/xorm"
	"windcontrol-go/persist/mysql"
)

var mysqlClient *xorm.Engine

func init()  {
	var err error
	if mysqlClient, err = mysql.NewClient(); err != nil {
		panic(err)
	}
}