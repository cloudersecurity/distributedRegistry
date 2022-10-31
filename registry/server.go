package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// 注册服务的逻辑
const (
	ServerPort  = ":3000"
	ServicesURL = "http://localhost" + ServerPort + "/services" // 启动注册服务后，通过这个URL查询有哪些服务
)

type registry struct {
	registrations []Registration // 存储已经注册的服务，因为是动态变化，要保证并发安全，所以加锁
	mutex         *sync.Mutex
}

// 注册服务的方法
func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.registrations = append(r.registrations, reg)
	// 注册服务的时候，向注册中心发送请求，把他依赖的服务请求过来？？
	err := r.sendRequiredServices(reg)
	// 根据本个变更patch信息，通知依赖这些服务的服务变更信息
	// 对于patch里新增的，通知依赖这些新增的服务的服务，他们所依赖的服务新增了
	// 对于patch里删除的，通知依赖这些删除的服务的服务，他们所依赖的服务删除了
	r.notify(patch{
		Added: []patchEntry{
			{
				Name: reg.ServiceName,
				URL:  reg.ServiceURL},
		},
	})
	return err
}

func (r *registry) notify(fullPatch patch) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, reg := range r.registrations {
		go func(reg Registration) {
			for _, reqService := range reg.RequiredServices {
				p := patch{Added: []patchEntry{}, Removed: []patchEntry{}}
				sendUpdate := false
				for _, added := range fullPatch.Added {
					if added.Name == reqService {
						p.Added = append(p.Added, added)
						sendUpdate = true
					}
				}
				for _, removed := range fullPatch.Removed {
					if removed.Name == reqService {
						p.Removed = append(p.Removed, removed)
						sendUpdate = true
					}
				}
				if sendUpdate {
					err := r.sendPatch(p, reg.ServiceUpdateURL)
					if err != nil {
						log.Println(err)
						return
					}
				}
			}
		}(reg)
	}
}

func (r *registry) sendRequiredServices(reg Registration) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	// 从已经注册的服务里边找reg这个新注册的服务所依赖的服务
	// 如果能找到就加到patch的added里边
	var p patch
	for _, serviceReg := range r.registrations {
		for _, reqService := range reg.RequiredServices {
			if serviceReg.ServiceName == reqService {
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceURL,
				})
			}
		}
	}
	//添加完patch之后，向通知URL将patch给被注册服务发送过去，通知它
	err := r.sendPatch(p, reg.ServiceUpdateURL)
	if err != nil {
		return err
	}
	return nil
}

func (r *registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return err
	}
	return nil
}

// 取消注册的方法
func (r *registry) Remove(url string) error {
	for i := range reg.registrations {
		if reg.registrations[i].ServiceURL == url {
			r.notify(patch{
				Removed: []patchEntry{
					{
						Name: r.registrations[i].ServiceName,
						URL:  r.registrations[i].ServiceURL,
					},
				},
			})
			reg.mutex.Lock()
			defer reg.mutex.Unlock()
			reg.registrations = append(reg.registrations[:i], reg.registrations[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Service at Url %s not found!", url)
}

// 检查服务心跳的方法
// 每隔一段时间向所有注册的服务发送请求，根据请求结果判断服务是否正常运行，以达到监控的效果
func (r *registry) heartBeat(freq time.Duration) {
	for {
		var wg sync.WaitGroup
		for _, reg := range r.registrations {
			wg.Add(1)
			go func(reg Registration) {
				defer wg.Done()
				success := true
				for attempts := 0; attempts < 3; attempts++ {
					resp, err := http.Get(reg.HeartBeatURL)
					if err != nil {
						log.Println(err)
					} else if resp.StatusCode == http.StatusOK {
						log.Printf("Heartbeat check passed for %v", reg.ServiceName)
						if !success {
							r.add(reg)
						}
						break
					}
					log.Printf("HeartBeat check failed for %v\n", reg.ServiceName)
					if success {
						success = false
						r.Remove(ServicesURL)
					}
					time.Sleep(time.Second)
				}
			}(reg)
			wg.Wait()
			time.Sleep(freq)
		}
	}
}

var once sync.Once

// 启动心跳检查协程
func SetupRegistryService() {
	once.Do(func() {
		go reg.heartBeat(3 * time.Second)
	})
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.Mutex),
}

type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("Request Received")
	switch req.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var r Registration
		err := decoder.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		log.Printf("serice: %v added with Url: %s\n", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		payload, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)
		log.Printf("Removing Service at URL: %s", url)
		err = reg.Remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
