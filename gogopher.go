package main

import (
	"fmt"
	"os"

	"gogopher/libgogo"
)

func main() {
	server, err := libgogo.ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf("Listening on %s:%d\n", server.Addr, server.Lport)
	err = server.Go()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
