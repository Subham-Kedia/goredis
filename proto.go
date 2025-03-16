package main

const (
	CommandSet  = "SET"
	CommandGet  = "GET"
	CommandQuit = "quit"
)

type Command any

type SetCommand struct {
	key, value string
}

type GetCommand struct {
	key string
}

type QuitCommand struct {
}
