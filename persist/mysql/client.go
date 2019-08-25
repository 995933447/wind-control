package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"windcontrol-go/config/persist"
	"xorm.io/core"
)

func NewClient() (*xorm.Engine, error) {
	client, err := xorm.NewEngine("mysql", persist.MysqlDsn)
	if err != nil {
		return nil, err
	}
	client.SetMaxIdleConns(10)
	client.SetMaxOpenConns(100)
	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, persist.MysqlPrefix)
	client.SetTableMapper(tbMapper)
	return client, err
}