package main

const (
	CommandSet = "SET"
	CommandGet = "GET"
)

type Command any

type SetCommand struct {
	key, value string
}

type GetCommand struct {
	key string
}
