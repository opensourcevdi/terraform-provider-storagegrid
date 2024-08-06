package main

import (
	"fmt"
	"log"
	"os"
	"terraform-provider-storagegrid/client"
)

func unwrap[T any](res T, err error) T {
	if err != nil {
		panic(err)
	}
	return res
}

func main() {
	c := client.Client{
		ApiUrl:    os.Getenv("apiUrl"),
		AccountId: os.Getenv("accountId"),
		Username:  os.Getenv("username"),
		Password:  os.Getenv("password"),
	}
	users := unwrap(c.GetUsers())
	log.Printf("%#v\n", users)
	for _, name := range os.Args[1:] {
		user := unwrap(c.GetUserByName(name))
		log.Printf("%#v\n", user)
		keys := unwrap(c.GetAccessKeys(user.Id))
		fmt.Printf("%#v\n", keys)
	}
}
