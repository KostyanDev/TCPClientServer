package model

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type User struct {
	name     string
	incoming chan *Message
	outgoing chan string
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	command  *Command
}

func (c *User)readInput(){
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		var comMap = map[string]{
			"/nick": {
				id:     NICK,
				client: c,
				args:   args,

			},
			"/quit": {
				id:     QUIT,
				client: c,
				args:   args,
			},
			"/members": {
				id:     MEMBERS,
				client: c,
				args:   args,
			},
		}

		k, ok := comMap[cmd]
		if !ok {
			fmt.Errorf("unknown command: %s", cmd)
		} else {
			c.command <- k
		}

	}
}
