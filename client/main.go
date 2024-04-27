package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const BUFF_SIZE = 2048

var input = make([]byte, BUFF_SIZE)

func handleError(err error) {
	if err == nil {
		return
	}
	fmt.Printf("error:%s\n", err.Error())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:<command> <port>")
		return
	}
	port := os.Args[1]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:"+port)
	handleError(err)
	tcpConn, err := net.DialTCP("tcp4", nil, tcpAddr)
	handleError(err)
	reader := bufio.NewReader(os.Stdin)
	var continued = true
	var inputStr string
	for continued {
		n, err := reader.Read(input)
		handleError(err)
		if n > 0 {
			k, _ := tcpConn.Write(input[:n])
			if k > 0 {
				inputStr = string(input[:k])
				fmt.Printf("Write:%s", inputStr)
				if inputStr == "exit\n" { //在比对时由于有个回车符，所以加上\n
					continued = false //也可以将inputStr = TrimRight(inputStr,"\n")
				}
			}
		}
	}
}
