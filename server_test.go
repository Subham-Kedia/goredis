package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/Subham-Kedia/goredis/client"
	"github.com/stretchr/testify/assert"
)

func TestServerWithMultipleClients(t *testing.T) {
	server := NewServer(Config{})
	go func() {
		log.Fatal(server.Start())
	}()

	time.Sleep(time.Second)

  nClients := 8
	wg := sync.WaitGroup{}
	wg.Add(nClients)
	for i := range nClients {
		go func(it int) {
			c, err := client.NewClient("localhost:5001")
			if err != nil {
				log.Fatal(err)
			}
			key := fmt.Sprintf("client_foo_%d", i)
			value := fmt.Sprintf("client_bar_%d", i)

			val, err := c.Set(context.TODO(), key, value)
			if err != nil {
				log.Fatal(err)
			}
      assert.Equal(t, val, "+OK\r\n")
			val, err = c.Get(context.TODO(), key)
			if err != nil {
				log.Fatal(err)
			}
      assert.Equal(t, value, val)
			c.Close()
			wg.Done()
		}(i)
	}
	wg.Wait()
}
