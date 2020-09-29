package model

import (
	"net"
	"strings"
	"time"
)

type Message struct {
	time    time.Time
	user    *User
	text    string
	command *Command
	conn 	net.Conn
	messages []string
	server 	*Server
}

func NewMessage(time time.Time, user *User, text string) *Message {
	return &Message{
		time: time,
		user: user,
		text: text,
	}
}


func (message *Message)ReadInput(msg *Message, users []*User) {

	tmp := "/"
	if strings.Contains(msg.text, tmp) == true {
			message.command.CommandRead(msg)
	} else {
		message.SendAll(msg.text, users)
	}
	//for {
		//msg, err := bufio.NewReader(message.conn).ReadString('\n')
		//if err != nil {
		//	return
		//}


		//msgString := string(msg.text)
		//msgString = strings.Trim(msgString, "\r\n")
		//
		//args := strings.Split(msg.text, " ")
		//cmd := strings.TrimSpace(args[0])

		//
		//message.command.CommandRead(cmd)

	//}
}

func (m *Message) SendAll(msg string, users []*User) {
	for _, user := range users {
		user.outgoing <- msg
	}
}
//func (msg *Message) String() string {
//	return fmt.Sprintf("%v [%v]: %v", msg.user.name, msg.time.Format("15:04:00"), msg.text)
//}
//
//func (c *Command) Parse(msg *Message) {
//	args := strings.TrimSpace(msg.text)
//	arrStr := strings.SplitN(args, " ", 2)
//
//	command := arrStr[0]
//	text := ""
//	if len(arrStr) == 2 {
//		text = arrStr[1]
//	}
//}


//var comMap = map[string]{
//	"/nick": {
//		id:     NICK,
//		client: c,
//		args:   args,
//	},
//	"/quit": {
//		id:     QUIT,
//		client: c,
//		args:   args,
//	},
//	"/members": {
//		id:     MEMBERS,
//		client: c,
//		args:   args,
//	},
//}
//
//k, ok := comMap[cmd]
//if !ok {
//	fmt.Errorf("unknown command: %s", cmd)
//} else {
//	c.command <- k
//}