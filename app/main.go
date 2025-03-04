package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		cmd := strings.Split(command, " ")
		if err == nil {
			switch cmd[0] {

			case "exit":
				os.Exit(0)

			case "echo":
				fmt.Printf("%s", strings.Join(cmd[1:], " "))

			default:
				fmt.Printf("%s: command not found\n", command[:len(command)-1])
			}
		}
	}
}
