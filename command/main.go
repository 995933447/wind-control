package main

import (
	"flag"
	"windcontrol-go/command/queue"
	"windcontrol-go/logger"
)

var engine string

func main()  {
	flag.StringVar(&engine,"engine", "", "请输入要启动的服务引擎,选项:queue队列服务引擎")
	flag.Parse()

	switch engine {
		case "queue":
			logger.DefaultLogger.Debug("Queue engine is running.", nil)
			queue.Run()
	}
}