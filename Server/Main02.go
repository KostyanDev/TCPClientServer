package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"reflect"
	"regexp"
	"strings"
)


type Client struct {
	name string
	id   int
	conn net.Conn
	commands chan<-command
}

// Hold the user list
var users = make(chan []string)
type Server struct {
	allClients       map[net.Conn]Client
	allNames        []string
	totalClients     int
	connectedClients int
	listener         net.Listener
}


var validName = regexp.MustCompile("^([A-Z]|[a-z]|[0-9]){4,12}$")

func getValidName(conn net.Conn) string {
	// Get name from client
	conn.Write([]byte("Please enter a new name : \n"))
	reader := bufio.NewReader(conn)
	name, _ := reader.ReadString('\n')
	name = strings.Trim(name, "\n")

	// Until it has correct format
	for !validName.MatchString(name) {
		conn.Write([]byte("Name are alphanumerical and of length in [4,12]\n"))
		conn.Write([]byte("Please enter a new name : \n"))
		name, _ = reader.ReadString('\n')
		name = strings.Trim(name, "\n")
	}
	return name
}

func disconnect(conn net.Conn, server *Server) {
	// Properly close the connection and delete the client from the server
	client := server.allClients[conn]
	name := client.name
	i := find(server.allNames, name)
	server.allNames = append(server.allNames[:i], server.allNames[i+1:]...)
	users <- server.allNames
	server.connectedClients--
	conn.Close()
	delete(server.allClients, conn)
	log.Printf("Client with name %s disconnected", name)
	log.Printf("There are %d connected clients", server.connectedClients)
}

func find(s interface{}, elem interface{}) int {
	// Return -1 if elem is in s, its index in s otherwise
	arrV := reflect.ValueOf(s)
	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			if arrV.Index(i).Interface() == elem {
				return i
			}
		}
	}
	return -1

}

func contains(s interface{}, elem interface{}) bool {
	// Return true if elem is in s, false otherwise
	arrV := reflect.ValueOf(s)
	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {

			// panics if slice element points to an unexported struct field
			if arrV.Index(i).Interface() == elem {
				return true
			}
		}
	}
	return false
}

func (c *Client)readInput(){
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		var comMap = map[string]command{
			"/nick": {
				id:     NICK,
				client: c,
				args:   args,
			},
			"/join": {
				id:     JOIN,
				client: c,
				args:   args,
			},
			"/rooms": {
				id:     ROOMS,
				client: c,
			},
			"/msg": {
				id:     MSG,
				client: c,
				args:   args,
			},
			"/quit": {
				id:     QUIT,
				client: c,
			},
			"/members": {
				id:     MEMBERS,
				client: c,
			},
		}

		k, ok := comMap[cmd]
		if !ok {
			c.err(fmt.Errorf("unknown command: %s", cmd))
		} else {
			c.commands <- k
		}
	}
}

func main() {
	// Initialization

	server := new(Server)
	server.allClients = make(map[net.Conn]Client)
	// Server will push new connections to it
	newConnections := make(chan net.Conn)
	// Clients that will be remove from allClients
	deadConnections := make(chan net.Conn)
	// Receives messages from connected clients
	messages := make(chan string)

	// Start TCP server

	var err error
	server.listener, err = net.Listen("tcp", ":6060")
	if err != nil {
		err = fmt.Errorf("error launching the server : %e", err)
		fmt.Println(err)
	}

	// Server accepts connections forever and pushes new ones to the channel
	go func() {
		for {
			conn, err := server.listener.Accept()
			if err != nil {
				fmt.Println(err)
			}
			newConnections <- conn
		}
	}()

	// Send the user list everytime it is modified
	go func() {
		// A user-list datagram will start with "##", which means
		// it is not a classic message since names are alphanumerical
		var b strings.Builder
		for {
			b.Reset()
			b.WriteString("##")
			userList := <-users
			for _, name := range userList {
				b.WriteString(name + ",")
			}
			b.WriteString("\n")
			userListMessage := b.String()
			messages <- userListMessage
		}
	}()

	for {
		select {

		// Continuously accept new clients
		case conn := <-newConnections:
			server.totalClients++
			log.Printf("Accepted new client with id %d", server.totalClients)

			// Read all incoming messages from this client and push them to the chan
			go func(conn net.Conn, server *Server) {
				conn.Write([]byte("Welcome to the server ! \n"))

				// Get a name in valid format
				name := getValidName(conn)
				// Get a name not used
				for contains(server.allNames, name) {
					conn.Write([]byte("Name already in use, please choose a new one"))
					name = getValidName(conn)
				}

				messages <- fmt.Sprintf("User %s joined the room !\n", name)
				conn.Write([]byte(fmt.Sprintf("Your name is now %s \n", name)))

				// Add client to server
				client := Client{name: name, id: server.totalClients}
				server.allClients[conn] = client
				server.allNames = append(server.allNames, name)
				users <- server.allNames
				server.connectedClients++
				reader := bufio.NewReader(conn)

				// Read all his incoming messages
				for {
					incoming, err := reader.ReadString('\n')
					if err != nil {
						break
					}
					messages <- fmt.Sprintf("%s > %s", client.name, incoming)
				}

				// If there was an error, we delete the client
				deadConnections <- conn
			}(conn, server)


		// Continuously read incoming messages and broadcast them
		case message := <-messages:
			log.Printf("New message : %s", message)

			// Send the message to all connected clients
			for conn := range server.allClients {
				//Send the message in a go routine
				go func(conn net.Conn, message string) {
					_, err := conn.Write([]byte(message))
					// If it doesn't work the connection is dead
					if err != nil {
						deadConnections <- conn
					}
				}(conn, message)
			}

		//Remove dead clients
		case conn := <-deadConnections:
			disconnect(conn, server)
		}
	}
}