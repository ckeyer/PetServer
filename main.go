package main

import (
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
	Run()
}

// Run 运行服务端，开始监听
func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Index)
	mux.HandleFunc("/upload", Receive)

	log.Notice("Http is running at ", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		panic(err)
	}
}
