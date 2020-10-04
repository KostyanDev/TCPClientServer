package model

import (
	"net"
	"strings"
)

type Message struct {

	user     *User
	text     string
	conn     net.Conn
	messages []string
	server   *Server
}

func NewMessage( user *User, text string, conn net.Conn) *Message {
	return &Message{
		user: user,
		text: text,
		conn: conn,
	}
}

func (message *Message) ReadInput(msg *Message, users []*User, user *User) {

	tmp := "/"
	if strings.Contains(msg.text, tmp) == true {
		CommandRead(msg)
	} else {
		message.SendAll(msg.text, users, user)
	}
}

func (m *Message) SendAll(msg string, users []*User, currentUser *User) {

	for _, user := range users {
		user.outgoing <- msg
	}
}
