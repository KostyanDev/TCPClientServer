package model

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const BUFFERSIZE = 1024

type File struct {
	name string
	conn net.Conn
}

func NewFile(name string, conn net.Conn) *File {
	return &File{
		name: name,
		conn: conn,
	}
}

func (f *File) listOfFile(user *User) {
	var files []string
	var (
		_, _, _, _ = runtime.Caller(0)
		basepath   = filepath.Dir("./")
	)

	root := basepath
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
	user.outgoing <- strings.Join(files[:], ", \n")
	user.outgoing <- "\n"
}

func (f *File) sendFiles(name string) {
	fmt.Println("A client has connected!")
	//defer connection.Close()
	//time.Sleep(20 * time.Second)
	file, err := os.Open("Test")
	fmt.Println("file is :", file)
	if err != nil {
		fmt.Println("file isn't open :", err)
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
	_, err = f.conn.Write([]byte(fileSize))
	if err != nil {
		log.Println("fileSize bad valye ")
	}
	_, err = f.conn.Write([]byte(fileName))
	if err != nil {
		log.Println("fileName bad valye ")
	}

	sendBuffer := make([]byte, BUFFERSIZE)

	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			log.Printf("Error is : %s", err)
			break
		}
		f.conn.Write(sendBuffer)
	}
	fmt.Println("File has been sent!")
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
