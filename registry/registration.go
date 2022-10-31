package registry

type Registration struct {
	ServiceName      ServiceName
	ServiceURL       string
	RequiredServices []ServiceName // 本服务依赖的其它服务，此例就是业务服务依赖日志服务
	ServiceUpdateURL string        // 注册服务通过这个URL来通知业务服务，它依赖的某服务可以用，此例就是注册服务通过这个URL告知业务服务，日志服务可用
	// 客户端的服务需要暴露这样的URL，来让服务注册中心通知它有服务更新了，从而实现服务发现
	HeartBeatURL string // 用于接收注册服务法送的心跳检测请求
}
type ServiceName string

const (
	LogService     = ServiceName("LogService")
	GradingService = ServiceName("GradingService")
)

// 变更的服务的条目
type patchEntry struct {
	Name ServiceName
	URL  string
}
type patch struct {
	Added   []patchEntry // 增加的条目
	Removed []patchEntry // 减少的条目
}
