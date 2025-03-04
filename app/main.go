package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err == nil {
			if command[:4] == "exit" {
				os.Exit(0)
			}
			fmt.Printf("%s: command not found\n", command[:len(command)-1])
		}
	}
}
