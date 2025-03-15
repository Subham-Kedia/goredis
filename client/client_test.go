package client

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(":5001")
	if err != nil {
		log.Fatal(err)
	}
	names := []string{"Anish", "Amit", "Abhishek"}
	for i, name := range names {
		err := client.Set(context.TODO(), fmt.Sprintf("name%d", i), name)
		if err != nil {
			log.Fatal(err)
		}
		val, err := client.Get(context.TODO(), fmt.Sprintf("name%d", i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(val)
	}
}

func TestMultipleClients(t *testing.T) {
	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)
  for i := range nClients {
		go func(it int) {
			c, err := NewClient("localhost:5001")
			if err != nil {
				log.Fatal(err)
			}
      defer c.Close()
			key := fmt.Sprintf("client_foo_%d", i)
			value := fmt.Sprintf("client_bar_%d", i)

			if err := c.Set(context.TODO(), key, value); err != nil {
				log.Fatal()
			}
			val, err := c.Get(context.TODO(), key)
			if err != nil {
				log.Fatal()
			}
			fmt.Printf("client %s got this val back", val)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
