package main

import (
	"fmt"
	"phonebook/_example/cmd/server/api"
)

// main
func main() {
	cli := api.NewServerCLI()

	if err := cli.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
