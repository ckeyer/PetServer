package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

var (
	url   = "http://127.0.0.1:8080/upload?test=true&debug=false"
	Token = []byte("123")
)

func main() {
	Upload(url, "a.txt")
}

func Upload(url, path string) (err error) {
	// Create buffer
	buf := new(bytes.Buffer) // caveat IMO dont use this for large files, \
	// create a tmpfile and assemble your multipart from there (not tested)
	w := multipart.NewWriter(buf)
	// defer w.Close()
	// Create file field
	fw, err := w.CreateFormFile("file", "test.txt") //这里的file很重要，必须和服务器端的FormFile一致
	if err != nil {
		fmt.Println("c")
		return err
	}

	fd, err := os.Open(path)
	if err != nil {
		fmt.Println("d")
		return err
	}
	filebs, _ := ioutil.ReadAll(fd)
	fmt.Printf("Fild is :%s\n", filebs)

	defer fd.Close()
	// Write file field from file to upload
	_, err = fw.Write(filebs)
	if err != nil {
		fmt.Println("e")
		return err
	}
	// Important if you do not close the multipart writer you will not have a
	// terminating boundry
	w.Close()
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		fmt.Println("f")
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	msha := HmacSha1(filebs, Token)
	fmt.Println(msha)
	req.Header.Set("X-Ckeyer-Sha1", msha)

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("g")
		return err
	}
	io.Copy(os.Stderr, res.Body) // Replace this with Status.Code check
	fmt.Println("h")
	return err
}

func HmacSha1(message, key []byte) string {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return fmt.Sprintf("%x", expectedMAC)
}

//客户端上传文件代码：
func Upload2() (err error) {
	// Create buffer
	buf := new(bytes.Buffer) // caveat IMO dont use this for large files, \
	// create a tmpfile and assemble your multipart from there (not tested)
	w := multipart.NewWriter(buf)
	// Create file field
	fw, err := w.CreateFormFile("file", "helloworld.go") //这里的file很重要，必须和服务器端的FormFile一致
	if err != nil {
		fmt.Println("c")
		return err
	}
	fd, err := os.Open("helloworld.go")
	if err != nil {
		fmt.Println("d")
		return err
	}
	defer fd.Close()
	// Write file field from file to upload
	_, err = io.Copy(fw, fd)
	if err != nil {
		fmt.Println("e")
		return err
	}
	// Important if you do not close the multipart writer you will not have a
	// terminating boundry
	w.Close()
	req, err := http.NewRequest("POST", "http://192.168.2.127/configure.go?portId=2", buf)
	if err != nil {
		fmt.Println("f")
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("g")
		return err
	}
	io.Copy(os.Stderr, res.Body) // Replace this with Status.Code check
	fmt.Println("h")
	return err
}
