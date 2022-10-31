package log

import (
	"bytes"
	"distributed/registry"
	"fmt"
	stdlog "log"
	"net/http"
)

// 为了更方便的启动Log Service，开发一个client
func SetClientLogger(serviceURL string, clientService registry.ServiceName) {
	stdlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	stdlog.SetFlags(0)
	stdlog.SetOutput(&clientLogger{url: serviceURL})
}

type clientLogger struct {
	url string
}

func (cl clientLogger) Write(data []byte) (int, error) {
	b := bytes.NewBuffer([]byte(data))
	resp, err := http.Post(cl.url+"/log", "text/plain", b)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Failed to send log message, Service responsed with %v", resp.StatusCode)
	}
	return len(data), nil
}
