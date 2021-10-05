package main

import (
	"fmt"
	"phonebook/cmd/client/api"
)

// main
func main() {
	cli := api.NewClientCLI()

	if err := cli.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
