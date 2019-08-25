package mysql

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"os"
	"strings"
	"util/filesystem"
	"util/timer"
)

const logPath  = "../../../storage/logs/mysql"

func ListenSql(engine *xorm.Engine) error {
	filename := fmt.Sprintf("%s/%s.log", strings.TrimRight(logPath, "/"), timer.NowDate("Y-m-d"))

	dir := filesystem.Dir(filename)
	if exist, err := filesystem.PathExists(dir); err != nil {
		return err
	} else if !exist {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}

	engine.ShowSQL(true)
	engine.SetLogger(xorm.NewSimpleLogger(file))
	return nil
}
