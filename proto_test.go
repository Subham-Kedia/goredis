package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"testing"

	"github.com/tidwall/resp"
)

func TestProtocol(t *testing.T) {
	msg := "*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n"
	rd := resp.NewReader(bytes.NewBufferString(msg))

	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Read: %s\n", v.Type())
		if v.Type() == resp.Array {
			for i, v := range v.Array() {
				fmt.Printf("  %d: %s\n", i, v.String())
			}
		}
	}
}
