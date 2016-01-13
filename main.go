package main

import (
	"fmt"
	"math"
	"net/http"
)

// init  初始化
func init() {
	log.Notice("PetServer Starting ...")
}

var (
	addr = ":8080"
)

func main() {
	err := LoadConf("v1.json")
	if err != nil {
		log.Panic(err)
	}

	err = Run()
	if err != nil {
		log.Panic(err)
	}

}

// Run 运行服务端，开始监听
func Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Index)
	mux.HandleFunc("/upload", Receive)
	log.Notice("Http is running at ", addr)
	err := http.ListenAndServe(addr, mux)
	return err
}

func FmtSize(size int64) string {
	var s []string = make([]string, 6)

	s[0] = "Byte"
	s[1] = "KB"
	s[2] = "MB"
	s[3] = "GB"
	s[4] = "TB"
	s[5] = "PB"
	for i, k := range s {
		v := math.Pow(1024, float64(i))
		if float64(size)/v < 1024 {
			return fmt.Sprintf("%.2f %s", float64(size)/float64(v), k)
		}
	}
	return ""
}
