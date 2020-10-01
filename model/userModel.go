package model

import (
	"bufio"
	"log"
	"net"
	"time"
)

type User struct {
	incoming chan *Message
	outgoing chan string
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	command  *Command
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
		msg := NewMessage(time.Now(), u, str, u.conn)
		u.incoming <- msg
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
