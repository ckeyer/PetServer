package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	// 文件根目录
	RootDir = "/data/pet/"

	// 用于文件上传认证
	Token = "123"
	// SecretKey string = os.Getenv(FILE_SERVER_SECRET_KEY)
)

type File struct {
	Name   string
	Path   string
	Force  bool
	Rename bool
}

// 请求处理
func Index(w http.ResponseWriter, req *http.Request) {
	defer func() {
		req.Body.Close()
	}()
	req.Header.Add("Access-Control-Allow-Origin", "*")
	err := req.ParseForm()
	if err != nil {
		http.Error(w, `{"error":"ParseFrom Failed"}`, http.StatusBadRequest)
		log.Error(err)
		return
	}

	log.Info("func Get:", req.Method)
	w.Write([]byte(`{"code":0}`))
	// switch strings.ToUpper(req.Method) {
	// case "POST":
	// 	Push(w, req)
	// 	return
	// case "GET":
	// default:
	// 	http.Error(w, `{"error":"Method Error"}`, http.StatusMethodNotAllowed)
	// 	return
	// }

	// path := RootDir + req.URL.Path
	// if !filter(path) {
	// 	http.Error(w, `{"error":"Not Exists"}`, http.StatusNotFound)
	// 	return
	// }
	// fif, err := os.Stat(path)
	// if err != nil || fif.IsDir() {
	// 	http.Error(w, `{"error":"Not Exists"}`, http.StatusNotFound)
	// 	return
	// }
	// f, err := os.OpenFile(path, os.O_RDONLY, 0444)
	// if err != nil {
	// 	http.Error(w, `{"error":"Not Exists"}`, http.StatusNotFound)
	// 	return
	// }
	// io.Copy(w, f)
}

// 上传文件
func Receive(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	err := req.ParseForm()
	if err != nil {
		http.Error(w, `{"code":-1,"error":"ParseForm Failed"}`, http.StatusBadRequest)
		return
	}

	if strings.ToUpper(req.Method) != "POST" {
		http.Error(w, `{"code":-1,"error":"Method Error"}`, http.StatusMethodNotAllowed)
		return
	}

	buf := new(bytes.Buffer)
	dir := RootDir

	// 接收文件
	for k, v := range req.Form {
		log.Debug(k, v)
	}

	f, fh, err := req.FormFile("file")
	if err != nil {
		http.Error(w, `{"error":"ParseFileForm  Failed"}`, http.StatusBadRequest)
		log.Error(err.Error())
		log.Debugf("%#v", f)
		log.Debugf("%#v", fh)
		return
	}

	//
	_, err = io.Copy(buf, f)
	if err != nil {
		log.Error(err.Error())
		return
	}

	bs := buf.Bytes()
	log.Debugf("%s", bs)
	hsha1 := req.Header.Get("X-Ckeyer-Sha1")
	hsha2 := HmacSha1(buf.Bytes(), Token)
	log.Debugf("Sha1 Mac: %s", hsha1)
	log.Debugf("Sha2 Mac: %s", hsha2)
	if hsha2 != hsha1 {
		http.Error(w, `{"error":"Auth failed"}`, http.StatusNotAcceptable)
		return
	}

	path := dir + fh.Filename
	if !filter(path) {
		http.Error(w, `{"error":"Not Exists"}`, http.StatusNotFound)
		return
	}
	log.Error(fh.Filename)

	// finfo, err := os.Stat(dir)
	// if err != nil {
	// 	os.MkdirAll(dir, 0644)
	// } else if !finfo.IsDir() {
	// 	if force {
	// 		os.Remove(dir)
	// 		os.MkdirAll(dir, 0644)
	// 	} else {
	// 		http.Error(w, `{"error":"Dir is a exists File"}`, http.StatusBadRequest)
	// 		return
	// 	}
	// }

	// if _, err = os.Stat(path); err == nil {
	// 	if force {
	// 		os.Remove(path)
	// 	} else {
	// 		http.Error(w, `{"error":"File Exists"}`, http.StatusBadRequest)
	// 		return
	// 	}
	// }

	newf, err := os.Create(path)
	if err != nil {
		http.Error(w, `{"error":"Create File Failed"}`, http.StatusBadRequest)
		return
	}
	defer newf.Close()
	_, err = io.Copy(newf, buf)
	if err != nil {
		http.Error(w, `{"error":"Upload File Failed"}`, http.StatusBadRequest)
		return
	}

	w.Write([]byte(`{"success":"ok"}`))
}

func filter(url string) bool {
	err_seps := []string{"..", "~", ".go", "--"}
	for _, sep := range err_seps {
		if strings.Count(url, sep) > 0 {
			return false
		}
	}
	return true
}
