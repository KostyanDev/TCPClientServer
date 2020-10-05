package model

import (
	"fmt"
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

func NewMessage(user *User, text string, conn net.Conn) *Message {
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

	var newUsers []*User
	for _, val := range users {

		newUsers = append(newUsers, val)
	}
	fmt.Println("newUsers", newUsers)
	//var newUsers []*User
	//for index,val := range users{
	//	if val == currentUser {
	//		m.conn.RemoteAddr()
	//		newUsers = append(newUsers[:index], newUsers[index + 1:]...)
	//	}
	//}
	//fmt.Println("old users,", users)
	//fmt.Println("new users,", newUsers)
	//for _, user := range users {
	//	user.outgoing <- msg
	//}
}
