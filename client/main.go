package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	serverURL := "http://localhost:8080"

	// 发送 GET 请求到服务器的 /blocks 路径
	resp, err := http.Get(serverURL + "/blocks")
	if err != nil {
		log.Fatal("Failed to send request:", err)
	}
	defer resp.Body.Close()

	// 读取服务器返回的响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read response:", err)
	}

	// 打印服务器返回的响应
	fmt.Println("Server response:", string(body))
}
