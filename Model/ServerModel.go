package Model

import "net"

type Server struct {
	allClients       map[net.Conn]Client
	allPseudo        []string
	totalClients     int
	connectedClients int
	listener         net.Listener
}
