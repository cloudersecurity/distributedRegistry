package grades

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// 同样需要调用service包里的start来启动学生业务  这里的register理解为注册路由地址
func RegisterHandlers() {
	handler := new(studentsHandler)
	http.Handle("/students", handler)  // 查学生集合
	http.Handle("/students/", handler) //students后面还可以加id，查单个学生
}

type studentsHandler struct{}

func (sh studentsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// /students   分割为2个子字符串，第一个是空串    查询
	// /students/{id} 分割为3个子字符串		   	 查询
	// /students/{id}/grades 分割为4个子字符串      新增
	pathSegments := strings.Split(req.URL.Path, "/")
	switch len(pathSegments) {
	case 2:
		sh.getAll(w, req)
	case 3:
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.getOne(w, req, id)
	case 4:
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.addGrade(w, req, id)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (sh studentsHandler) getAll(w http.ResponseWriter, req *http.Request) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()
	data, err := sh.toJson(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (sh studentsHandler) toJson(obj interface{}) ([]byte, error) {
	var bf bytes.Buffer
	enc := json.NewEncoder(&bf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}

func (sh studentsHandler) getOne(w http.ResponseWriter, req *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()
	student, err := students.GetById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}
	data, err := sh.toJson(student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to serialize student data, err: %s", err.Error())
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}
func (sh studentsHandler) addGrade(w http.ResponseWriter, req *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()
	student, err := students.GetById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
	}
	var g Grade
	dec := json.NewDecoder(req.Body)
	err = dec.Decode(&g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	student.Grades = append(student.Grades, g)
	w.WriteHeader(http.StatusCreated)
	data, err := sh.toJson(g)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}
