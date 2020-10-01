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

	conn, err := net.Dial("tcp", ":8083")
	if err != nil {
		fmt.Println(err)
	}

	go DownloadFile(conn)
	//go Read(conn)
	//go Write(conn)

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

//func DownloadFile(conn io.ReadCloser) error {
//	//defer conn.Close()
//	reader := bufio.NewReader(conn)
//	// читаем первую строку - это будет название файла
//	line, err := reader.ReadString('\n')
//	if err != nil {
//		return err
//	}
//	line = strings.TrimSpace(line)
//	// создаём файл с заданным именем в текущей директории
//	file, err := os.Create(line)
//	if err != nil {
//		return err
//	}
//	// и копируем в него всё что дальше приходит от клиента
//	_, err = io.Copy(file, conn)
//	return err
//}
//
func DownloadFile(conn net.Conn) {
	defer conn.Close() // не надо ставить - это означает отключение связи клиента с сервером
	//time.Sleep(2 * time.Second)
	for {

		fmt.Println("Connected to server, start receiving the file name and file size")
		bufferFileName := make([]byte, 64)
		bufferFileSize := make([]byte, 10)
		conn.Read(bufferFileSize)
		fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

		conn.Read(bufferFileName)

		fileName := strings.Trim(string(bufferFileName), ":")

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
		fmt.Println("Received file completely!\n")
	}

}
