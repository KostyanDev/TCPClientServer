package model

import (
	"log"
	"net"
)

type Server struct {
	users      []*User
	allClients map[net.Conn]User
	message    *Message
	incoming   chan *Message
	join       chan *User
	leave      chan *User
}

func StartServer() *Server {
	server := &Server{
		users:    make([]*User, 0),
		incoming: make(chan *Message),
	}
	server.Listen()
	return server
}
func (s *Server) Listen() {
	go func() {
		for {
			select {
			case msg := <-s.incoming:
				//s.message.command.SendAll(msg, s.users)
				//s.SendAll(msg)
				s.message.ReadInput(msg,s.users)
				//fmt.Println("here a message", msg)
			case user := <-s.join:
				s.Join(user)
			case user := <-s.leave:
				s.Leave(user)
			}
		}
	}()
}

func (s *Server) Join(user *User) {
	s.users = append(s.users, user)
	log.Printf("New user (%v) joined in lobby.\n", user.conn.RemoteAddr().String())
	user.outgoing <- "Welcome to the Chat! Type \"/help\" to get a list commands.\n"

	// Это получает все сообщения с каждого юзера и перемещает в канал сервера.
	go func() {
		for message := range user.incoming {
			s.incoming <- message
		}
		s.leave <- user
	}()
	//-----
}

func (s *Server) Leave(user *User) {
	for index, val := range s.users {
		if user == val {
			s.users = append(s.users[:index], s.users[index:]...)
			break
		}
	}
	close(user.outgoing)
	log.Printf("Closed %s's outgoing channel.\n", user.name)
}

// Почитай про каналы и горутины.
func (s *Server) SendAll2(msg *Message) {
	for _, user := range s.users {
		user.outgoing <- msg.text
	}
}
//func (s *Server) SendMessage(msg *Message) {
//	s.message = append(s.message, msg)
//
//	for _, user := range s.users {
//		user.outgoing <- msg
//	}
//}//type Server struct {
////	users            []*User
////	messege          *Message
////	allClients       map[net.Conn]User
////	incoming         chan *Message
////	leave            chan *User
////	allPseudo        []string
////	totalClients     int
////	connectedClients int
////	listener         net.Listener
////}