package main

import (
	"bytes"
	"io"
	"log"

	"github.com/tidwall/resp"
)

const (
	CommandSet = "set"
	CommandGet = "get"
)

type Command interface {
}

type SetCommand struct {
	key, value string
}

type GetCommand struct {
	key string
}

func parseMessage(raw string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))
	var cmd Command
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if v.Type() == resp.Array {
			for i, c := range v.Array() {
				switch c.String() {
				case CommandSet:
					cmd = &SetCommand{
						key:   v.Array()[i+1].String(),
						value: v.Array()[i+2].String(),
					}
				case CommandGet:
					cmd = &GetCommand{
						key: v.Array()[i+1].String(),
					}
				default:
					return cmd, nil
				}
			}
		}
	}
	return cmd, nil
}
