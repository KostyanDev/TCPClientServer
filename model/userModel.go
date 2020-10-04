package model

import (
	"bufio"
	"log"
	"net"
)

type User struct {
	incoming chan *Message
	outgoing chan string
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func NewUser(conn net.Conn) *User {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	u := &User{
		incoming: make(chan *Message, 3),
		outgoing: make(chan string, 3),
		conn:     conn,
		reader:   reader,
		writer:   writer,
	}
	u.Listen()
	return u
}

func (u *User) Listen() {
	go u.Read()
	go u.Write()
}

func (u *User) Read() {
	for {
		str, err := u.reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			break
		}
		msg := NewMessage(u, str, u.conn)
		u.incoming <- msg

		//tmp := "/download"
		//if strings.Contains(str,tmp) == true{
		//	DownloadFile(u.conn)
		//}
	}
	close(u.incoming)
	log.Printf("Closed s incoming channel read thread")
}

func (u *User) Write() {
	for str := range u.outgoing {
		if _, err := u.writer.WriteString(str); err != nil {
			log.Println(err)
			break
		}

		if err := u.writer.Flush(); err != nil {
			log.Println(err)
			break
		}
	}
	log.Printf("Closed write thread")
}

//func DownloadFile(conn net.Conn) {
//	//defer conn.Close() // не надо ставить - это означает отключение связи клиента с сервером
//	//time.Sleep(2 * time.Second)
//
//	fmt.Println("Connected to server, start receiving the file name and file size")
//	bufferFileName := make([]byte, 64)
//	bufferFileSize := make([]byte, 10)
//	_, err := conn.Read(bufferFileSize)
//	if err != nil {
//		log.Println("read bufferFileSize bad ")
//	}
//	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
//
//	_, err = conn.Read(bufferFileName)
//	if err != nil {
//		log.Println("read 2 bufferFileSize bad ")
//	}
//
//	fileName := strings.Trim(string(bufferFileName), ":")
//
//	newFile, err := os.Create(fileName)
//	if err != nil {
//		log.Println("Error 180")
//		return
//	}
//	defer newFile.Close()
//	var receivedBytes int64
//
//	for {
//		if (fileSize - receivedBytes) < BUFFERSIZE {
//			io.CopyN(newFile, conn, (fileSize - receivedBytes))
//			conn.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
//			break
//		}
//		io.CopyN(newFile, conn, BUFFERSIZE)
//		receivedBytes += BUFFERSIZE
//	}
//	fmt.Println("Received file completely!")
//}