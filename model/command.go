package model

import (
	"log"
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
	//	command  *Command
	//	name     *User
	//	server   *Server
	file *File
	//	fileName *Message
	//	messages []string
	//}
}

func NewCommand(file *File) *Command {
	return &Command{
		file: file,
	}
}

func CommandRead(msg *Message) {

	args := strings.TrimSpace(msg.text)
	arrStr := strings.SplitN(args, " ", 2)

	command := arrStr[0]
	text := ""
	if len(arrStr) == 2 {
		text = arrStr[1]
	}
	file := newFile(text)
	switch command {
	case "/changeName":
		//ChangeName(c.name, text, c.server.allClients)
	case "/help":
		Help(msg.user)
	case "/list":
		//c.file.ListOfFile()
		file.ListOfFile(msg.user)

	case "/sendFile":

	case "/download":
		file.SendFiles(msg.conn, text)
	default:
		SendToUser(text, msg.user)
	}

}

func quit() {

}
func Help(user *User) {
	user.outgoing <- "\n\tCommands:\n"
	user.outgoing <- "\t/help - lists all commands.\n"
	user.outgoing <- "\t/name foo - changes your name to foo.\n"
	user.outgoing <- "\t/list - list of all file in folder.\n"
	user.outgoing <- "\t/download file - download with name file.\n"
	log.Printf("requested help.")
}

//func (c *Command) SendAll(msg *Message, users []*User) {
//	c.messages = append(c.messages, msg.text)
//
//	for _, user := range users  {
//		user.outgoing <- msg.text
//	}
//}

//func (s *Server) SendAll(msg *Message) {
//	for _, user := range s.users {
//		user.outgoing <- msg.text
//	}
//}
func SendToUser(msg string, user *User) {
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
