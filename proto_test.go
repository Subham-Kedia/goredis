package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/resp"
)

func TestProtocol(t *testing.T) {
	msg := "*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n"
	rd := resp.NewReader(bytes.NewBufferString(msg))
  
  expected := []string{"SET", "mykey", "myvalue"}

	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if v.Type() == resp.Array {
			for i, v := range v.Array() {
        assert.Equal(t, expected[i], v.String())
			}
		}
	}
}
