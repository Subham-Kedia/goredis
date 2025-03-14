package client

import (
	"context"
	"fmt"
	"log"
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
