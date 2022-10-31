package main

import (
	"context"
	"distributed/grades"
	"distributed/log"
	"distributed/registry"
	"distributed/service"
	"fmt"
	stlog "log"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)
	r := registry.Registration{
		ServiceName:      registry.GradingService,
		ServiceURL:       serviceAddress,
		RequiredServices: []registry.ServiceName{registry.LogService},
		ServiceUpdateURL: serviceAddress + "/services",
		HeartBeatURL:     serviceAddress + "/heart",
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		grades.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatal(err)
	}
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	logProvider, err := registry.GetProvider(registry.LogService)
	if err != nil {
		fmt.Printf("Failed to Get Provider with error: %v", err)
	} else if err == nil {
		fmt.Printf("logging service found at %s\n", logProvider)
		log.SetClientLogger(logProvider, r.ServiceName)
	}
	<-ctx.Done()
	fmt.Println("Shutting Down Grading Service!")
}
