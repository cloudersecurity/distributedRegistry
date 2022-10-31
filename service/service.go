package service

// 本包集中化启动服务（通用的），可以用来注册service，也可以启动service；并且给用户一个选项，可以随时停止service
import (
	"context"
	"distributed/registry"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, host, port string, reg registry.Registration,
	registerHandlersFunc func()) (context.Context, error) {
	registerHandlersFunc() // 注册路由地址的函数
	// 启动日志服务
	ctx = startService(ctx, reg.ServiceName, host, port)
	// 将日志服务注册到注册服务里，所以测试的时候先手动启动注册服务，再手动启动日志服务
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}
func startService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	var srv http.Server
	srv.Addr = host + ":" + port
	go func() {
		// 如果服务启动失败，错误就会被打印到日志里
		log.Println(srv.ListenAndServe())
		log.Println("failed to start")
		err := registry.ShutDownService(fmt.Sprintf("http://%s:%s", host, port)) // 删除日志注册服务
		if err != nil {
			log.Println(err) // 这里不return了，往下走cancel掉
		}
		cancel() //一旦发生错误就调用上下文的取消方法，告诉其它goroutin也停止，cancel是上下文包中用来给协程通信的方法
	}()

	go func() {
		fmt.Printf("%v started, Press any key to stop. \n", serviceName)
		var s string
		fmt.Scanln(&s)                                                           // 不按键的时候会等待在这里，如果按了，就会往下走，关闭服务
		err := registry.ShutDownService(fmt.Sprintf("http://%s:%s", host, port)) // 删除日志注册服务
		if err != nil {
			log.Printf("Failed to remove service, err%v\n", err)
		}
		err = srv.Shutdown(ctx)
		if err != nil {
			fmt.Printf("Failed to shut down Service, err: %v\n", err)
		}
		cancel()
	}()
	return ctx
}
