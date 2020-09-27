package model

import (
	"log"
	"net"
	"os/user"
	"strings"
)

type Server struct {
	users            []*User
	messege          *Message
	allClients       map[net.Conn]User
	incoming         chan *Message
	leave            chan *User
	allPseudo        []string
	totalClients     int
	connectedClients int
	listener         net.Listener
}

func startServer() *Server {
	server := &Server{
		users: make([]*User, 0),
	}
	server.Listen()
	return server
}
func (s *Server) Listen() {
	go func() {
		for {
			select {
			case msg := <-s.incoming:
				s.Parse(msg)
			case user := <-l.join:
				l.Join(user)
			case user := <-l.leave:
				l.Leave(user)
			}
		}
	}()
}

func (s *Server) SendMessage(message string) {
	s.SendAll(message)
}

func (s *Server) SendAll(msg string) {
	s.messege = append(s.messege, msg)

	for _, user := range s.users {
		user.outgoing <- msg
	}
}
