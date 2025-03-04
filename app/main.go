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
	builtinCommands := map[string]bool{"exit": true, "echo": true, "type": true}

	for {
		fmt.Fprint(os.Stdout, "$ ")

		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		cmd := strings.Split(strings.TrimSpace(command), " ")

		if err == nil {
			switch cmd[0] {

			case "exit":
				os.Exit(0)

			case "echo":
				fmt.Printf("%s\n", strings.Join(cmd[1:], " "))

			case "type":
				if builtinCommands[cmd[1]] {
					fmt.Printf("%s is a shell builtin\n", cmd[1])
				} else {
					fmt.Printf("%s: not found\n", cmd[1])
				}

			default:
				fmt.Printf("%s: command not found\n", command[:len(command)-1])
			}
		}
	}
}
