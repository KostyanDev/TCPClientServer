package model

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

const BUFFERSIZE = 1024

type File struct {
	name string
	size float32
}

func (f *File) ListOfFile() []string {
	var files []string
	root := "classic root "
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return err
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		f.name = file
	}
	return files
}
func (f *File) SendFiles(connection net.Conn, name string) {
	fmt.Println("A client has connected!")
	defer connection.Close()
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
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File has been sent, closing connection!")
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
