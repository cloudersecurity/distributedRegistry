package main

import (
	"context"
	"distributed/log"
	"distributed/registry"
	"distributed/service"
	"fmt"
	stlog "log"
)

func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "4000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)
	r := registry.Registration{
		ServiceName:      registry.LogService,
		ServiceURL:       fmt.Sprintf("http://%s:%s", host, port),
		RequiredServices: make([]registry.ServiceName, 0),
		ServiceUpdateURL: serviceAddress + "/services",
		HeartBeatURL:     serviceAddress + "/heart",
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		log.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatalln(err) // 如果服务启动失败，用标准库的log写日志，因为我们创建的log服务可能还没启动成功
	}
	<-ctx.Done() // Done是一个退出信号的chan，在本项目中，有2出cancel会触发done信号，一个是启动服务失败，另一个是手动停止服务
	fmt.Println("Shutting Down Log Service")
}
