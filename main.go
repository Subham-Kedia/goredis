package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Subham-Kedia/goredis/client"
)

func main() {
	// asynchronously starting the server
	server := NewServer(Config{ListenAddress: ":5001"})
	go func() {
		log.Fatal(server.Start())
	}()

	// waiting for the server to start
	time.Sleep(time.Second * 2)

	// creating a client and setting and getting 10 keys
	client := client.NewClient(":5001")
	for i := range 10 {
		err := client.Set(context.TODO(), fmt.Sprintf("foo%d", i), fmt.Sprintf("bar%d", i))
		if err != nil {
			log.Fatal(err)
		}
		// time.Sleep(time.Second)
		_, err = client.Get(context.TODO(), fmt.Sprintf("foo%d", i))
		if err != nil {
			log.Fatal(err)
		}
		// time.Sleep(time.Second)
	}
	// waiting for the client to finish and server to process the requests
	time.Sleep(time.Second * 2)

	// printing the data stored in the server
	fmt.Println(server.kv.data)
}
