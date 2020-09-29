package model

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type commandID int

//const (
//	NICK         commandID = 1
//	QUIT         commandID = 2
//	MEMBERS      commandID = 3
//	HELP         commandID = 4
//	ListOfFile   commandID = 5
//	SendFile     commandID = 6
//	DownloadFile commandID = 7
//)

type Command struct {
	command  *Command
	name     *User
	server   *Server
	file     *File
	fileName *Message
	messages []string
}

func (c *Command) CommandRead(msg *Message) {

	args := strings.TrimSpace(msg.text)
	arrStr := strings.SplitN(args, " ", 2)

	command := arrStr[0]
	text := ""
	if len(arrStr) == 2 {
		text = arrStr[1]
	}

	switch command {
	case "/changeName":
		ChangeName(c.name, text, c.server.allClients)
	case "/help":
		Help(c.name)
	case "/list":
		c.file.ListOfFile()
	case "/sendFile":

	case "/download":
		c.file.SendFiles(c.name.conn, "test.rtf")
	default:
		c.SendToUser(command, c.name)
	}

}

func ChangeName(user *User, name string, arraOfClients map[net.Conn]User) {

	for _, val := range arraOfClients {
		if name == val.name {
			fmt.Sprintf("nickname %s is taken", name)
		}
	}
	log.Printf("%s changed their name to %s", user.name, name)
	user.name = name
}
func quit() {

}
func Help(user *User) {
	user.outgoing <- "\n\tCommands:\n"
	user.outgoing <- "\t/help - lists all commands.\n"
	user.outgoing <- "\t/name foo - changes your name to foo.\n"
	user.outgoing <- "\t/list - list of all file in folder.\n\n"
	user.outgoing <- "\t/download - download file with name.\n\n"
	log.Printf("%s requested help.", user.name)
}
func (c *Command) SendAll(msg *Message, users []*User) {
	c.messages = append(c.messages, msg.text)

	for _, user := range users  {
		user.outgoing <- msg.text
	}
}

func (s *Server) SendAll(msg *Message) {
	for _, user := range s.users {
		user.outgoing <- msg.text
	}
}
func (c *Command) SendToUser(msg string, user *User){
	user.outgoing <- "\t/somesing gone wrong.\n\n"
}


//var comMap = map[string]{
//	"/changeName": {
//		ChangeName(c.client, c.name.name, c.server.allClients),
//	},
//	"/help":{
//		Help(c.client),
//	},
//	"/list":{
//		c.file.ListOfFile(),
//	},
//	"/download":{
//		c.file.SendFiles(c.client.conn, c.fileName.text),
//	},
//}
//k, ok := comMap[comm]
//if !ok {
//	fmt.Errorf("unknown command: %s", comm)
//} else {
//	return k
//}