package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

const BUFFERSIZE = 1024
func main() {
	wg.Add(1)

	conn, err := net.Dial("tcp", ":8082")
	if err != nil {
		fmt.Println(err)
	}


	DownloadFile(conn)
	go Read(conn)
	go Write(conn)

	wg.Wait()
}

func Read(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Disconnected from the server.")
			wg.Done()
			return
		}
		fmt.Print(str)
	}
}

func Write(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)

	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, err = writer.WriteString(str)
		if err != nil {
			os.Exit(2)
		}

		if err = writer.Flush(); err != nil {
			fmt.Println(err)
		}
	}
}

func DownloadFile(conn net.Conn){
	//defer conn.Close()
	fmt.Println("Connected to server, start receiving the file name and file size")
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)


	conn.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	conn.Read(bufferFileName)
	//fileName := strings.Trim(string(bufferFileName), ":")
	fileName := ("test.rtf")

	newFile, err := os.Create(fileName)

	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, conn, (fileSize - receivedBytes))
			conn.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, conn, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Println("Received file completely!")
}
