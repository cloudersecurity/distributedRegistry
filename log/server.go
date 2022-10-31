package log

// 日志是一个web服务，可以接收post请求，把请求内容写入日志文件里
// 本文件为日志服务后端逻辑
import (
	"io/ioutil"
	stlog "log"
	"net/http"
	"os"
)

var log *stlog.Logger

// 文件路径
type filelog string

func (fl filelog) Write(data []byte) (int, error) {
	// 打开文件，os.O_CREATE代表文件不存在就创建一个,打开的时候只写，且追加写
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)
}

// 根据目标文件地址创建logger，定义日志格式
func Run(destination string) {
	log = stlog.New(filelog(destination), "go: ", stlog.LstdFlags)
}

// 处理请求
func RegisterHandlers() {
	http.HandleFunc("/log", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			msg, err := ioutil.ReadAll(req.Body)
			if err != nil || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// 如果没有错误，把请求内容写入日志文件
			write(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}
func write(msg string) {
	log.Printf("%v\n", msg)
}
