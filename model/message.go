package model

import (
	"fmt"
	"strings"
	"time"
)

type Message struct {
	time    time.Time
	user    *User
	text    string
	command *Command
}

func NewMessage(time time.Time, user *User, text string) *Message {
	return &Message{
		time: time,
		user: user,
		text: text,
	}
}

func (msg *Message) String() string {
	return fmt.Sprintf("%v [%v]: %v", msg.user.name, msg.time.Format("15:04:00"), msg.text)
}

func (c *Command) Parse(msg *Message) {
	args := strings.TrimSpace(msg.text)
	arrStr := strings.SplitN(args, " ", 2)

	command := arrStr[0]
	text := ""
	if len(arrStr) == 2 {
		text = arrStr[1]
	}

	switch command {
	case "/help":
		c.Parse()
	case "/list":
		s.ListChatRooms(msg.user)
	case "/create":
		s.CreateChatRoom(msg.user, text)
	case "/join":
		l.JoinToChatRoom(msg.user, text)
	case "/leave":
		l.LeaveChatRoom(msg.user)
	case "/name":
		l.ChangeName(msg.user, text)
	case "/quit":
		msg.user.Quit()
	case "/del":
		l.DeleteChatRoom(msg.user, text)
	default:
		l.SendMessage(msg)
	}
}
