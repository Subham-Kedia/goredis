package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Subham-Kedia/goredis/client"
)

func main() {
	server := NewServer(Config{ListenAddress: ":5001"})
	// goroutine to start the server
	go func() {
		log.Fatal(server.Start())
	}()

	time.Sleep(time.Second)

	// 10 client requesting for setting values
	for i := range 10 {
		client := client.NewClient(":5001")
		err := client.Set(context.TODO(), fmt.Sprintf("foo%d", i), fmt.Sprintf("bar%d", i))
		if err != nil {
			log.Fatal(err)
		}
		_, err = client.Get(context.TODO(), fmt.Sprintf("foo%d", i))
		if err != nil {
			log.Fatal(err)
		}
	}

	time.Sleep(time.Second)
	fmt.Println(server.kv.data)
}
