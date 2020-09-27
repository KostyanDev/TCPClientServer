package model

import (
	"fmt"
	"log"
	"net"
)

type commandID int

const (
	NICK         commandID = iota
	QUIT         commandID = iota
	MEMBERS      commandID = iota
	HELP         commandID = iota
	ListOfFile   commandID = iota
	SendFile     commandID = iota
	DownloadFile commandID = iota
)

type Command struct {
	id       commandID
	client   *User
	name     *User
	server   *Server
	file     *File
	fileName *Message
}

func (c *Command) CommandReady(comm commandID) {
	switch comm {
	case NICK:
		ChangeName(c.client, c.name.name, c.server.allClients)
	case HELP:
		Help(c.client)
	case ListOfFile:
		c.file.ListOfFile()
	case SendFile:

	case DownloadFile:
		c.file.SendFiles(c.client.conn, c.fileName.text)
	}

}
func Parse() {

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
	user.outgoing <- "\t/list - lists all chat room.\n"
	user.outgoing <- "\t/create foo - creates a chat room with name \"foo\".\n"
	user.outgoing <- "\t/del foo - deletes a chat room.\n"
	user.outgoing <- "\t/join foo - joins a chat room named foo.\n"
	user.outgoing <- "\t/leave - leaves the current chat room.\n"
	user.outgoing <- "\t/name foo - changes your name to foo.\n"
	user.outgoing <- "\t/quit - quits the program.\n\n"
	log.Printf("%s requested help.", user.name)
}
