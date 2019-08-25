package createserverjs

import (
	"fmt"
	"os"
	"strings"
	filesystem "util/filesystem"
	"windcontrol-go/config/types"
	"windcontrol-go/logger"
	"windcontrol-go/queue/service/createserverjs/config"
)

func JobHandle(task types.Task) error {
	logger.DefaultLogger.Debug("start server.js task\n", nil)
	dir := strings.TrimRight(config.JsPath, "/") + "/"
	exist, err := filesystem.PathExists(dir)
	if err != nil {
		return err
	}
	if !exist {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	addresses := task.(map[string]map[string]string)
	for filename := range addresses {
		file, err := os.OpenFile(dir + filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return err
		}

		fileContent := ""
		for serverType, address := range addresses[filename] {
			if address == "" {
				logger.DefaultLogger.Emergency(serverType + " has not avaliabled address", nil)
			}
			fmt.Println(strings.TrimRight(address, "/") + "/", address)
			fileContent = fileContent + fmt.Sprintf("window.%s = %q\n", serverType, strings.TrimRight(address, "/") + "/")
		}

		logger.DefaultLogger.Debug("updating server js.", nil)
		_, err = file.Write([]byte(fileContent))
		if err != nil {
			return err
		}
	}

	return nil
}
