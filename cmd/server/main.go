package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

const BUFFERSIZE = 1024

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//server := model.StartServer()

	listener, err := net.Listen("tcp", ":8085")
	if err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}
	defer listener.Close()
	log.Println("Listening on localhost:8081.")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		//server.Join(model.NewUser(conn))
		SendFiles(conn, "Test")

	}
}
func SendFiles(connection net.Conn, name string) {
	fmt.Println("A client has connected!")
	//defer connection.Close()
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	fmt.Println("Sending filename and filesize!")
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!\n")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File has been sent! \n")
	return
}
func downloadFile(name string) {

}

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}
