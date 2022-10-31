package registry

// 此包用于注册服务的客户端,发post请求给注册服务
import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
)

func RegisterService(r Registration) error {
	heartBeatURL, err := url.Parse(r.HeartBeatURL)
	if err != nil {
		return err
	}
	http.HandleFunc(heartBeatURL.Path, func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	serviceUpdateURL, err := url.Parse(r.ServiceUpdateURL)
	if err != nil {
		return err
	}
	http.Handle(serviceUpdateURL.Path, &serviceUpdateHandler{})
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	err = encoder.Encode(r)
	if err != nil {
		return err
	}
	resp, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to register service, ServiceName: %s with status code: %v", r.ServiceName, resp.StatusCode)
	}
	return nil
}

type serviceUpdateHandler struct{}

func (sup serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(req.Body)
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("Updated recevied %v\n", p)
	prov.Update(p)
}

// 再写一个取消注册服务的函数，充当客户端
func ShutDownService(url string) error {
	fmt.Println("t1")
	req, err := http.NewRequest(http.MethodDelete, ServicesURL, bytes.NewBuffer([]byte(url)))
	fmt.Println("t2")
	if err != nil {
		return err
	}
	fmt.Println("t3")
	req.Header.Add("Content-Type", "text/plain")
	fmt.Println("t4")
	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("%#v", resp)
	fmt.Println("t5")

	if err != nil {
		return err
	}
	fmt.Println("t6")

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to deregister service with code %v", resp.StatusCode)
	}
	return nil
}

// 每个客户端的服务都有自己所依赖的服务，客户端需要有个地方存储他们所请求的服务
// 给业务服务提供服务的服务，比如这里的日志服务给学生成绩服务提供服务
type providers struct {
	services map[ServiceName][]string // 之所以是字符串切片，因为日志可能有多个URL来提供服务，本例中只有1个
	mutex    *sync.RWMutex
}

func (p *providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// 根据传进来的patch来更新provider
	// 首先是新增部分
	for _, patchEntry := range pat.Added {
		//如果要新增的这个服务当前还不存在,就在map里创建
		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = make([]string, 0)
		}
		// 如果存在，就在后面添加上URL
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL)
	}
	// 下面是删除的情况
	for _, patchEntry := range pat.Removed {
		if provierURLs, ok := p.services[patchEntry.Name]; ok {
			for i := range provierURLs {
				if provierURLs[i] == patchEntry.URL {
					p.services[patchEntry.Name] = append(provierURLs[:i], provierURLs[i+1:]...)
				}
			}
		}
	}
}

// 通过服务的名称，找到它所依赖服务的URL
func (p providers) get(name ServiceName) (string, error) {
	providers, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("No Poviders available for service %v", name)
	}
	if len(providers) > 0 {
		return providers[0], nil
	}
	return "", nil
}

// 将get上面套一个公有方法
func GetProvider(name ServiceName) (string, error) {
	return prov.get(name)
}

// 创建变量用来存储接收到的patch
var prov = providers{services: make(map[ServiceName][]string, 0), mutex: new(sync.RWMutex)}
