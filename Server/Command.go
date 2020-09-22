package main

type commandID int

const (
	NICK commandID = iota
	JOIN
	ROOMS
	MSG
	QUIT
	MEMBERS
)

type command struct{
	id commandID
	client *client
	args []string
}