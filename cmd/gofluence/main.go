package main

import (
	"fmt"
	"os"

	"github.com/beto20/gofluence/command"
)

func main() {
	if err := command.Root(os.Args[0:]); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
