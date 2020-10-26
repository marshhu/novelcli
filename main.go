package main

import (
	"fmt"
	"os"

	"github.com/marshhu/novelcli/cmd"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
	}()
	cmd.Execute()
}
