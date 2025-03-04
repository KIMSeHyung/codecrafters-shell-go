package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	builtinCommands := map[string]bool{"exit": true, "echo": true, "type": true, "pwd": true}
	path := strings.Split(os.Getenv("PATH"), ":")

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
					found := false
					for _, dir := range path {
						fullPath := filepath.Join(dir, cmd[1])
						if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
							fmt.Printf("%s is %s\n", cmd[1], fullPath)
							found = true
							break
						}
					}
					if !found {
						fmt.Printf("%s: not found\n", cmd[1])
					}
				}

			case "pwd":
				dir, _ := os.Getwd()
				fmt.Println(dir)

			default:
				if len(cmd) > 1 {
					found := false
					for _, dir := range path {
						fullPath := filepath.Join(dir, cmd[0])
						if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
							c := exec.Command(cmd[0], cmd[1:]...)
							c.Stdout = os.Stdout
							c.Stderr = os.Stderr
							c.Run()
							found = true
							break
						}
					}
					if !found {
						fmt.Printf("%s: not found\n", cmd[0])
					}
				} else {
					fmt.Printf("%s: command not found\n", command[:len(command)-1])
				}
			}
		}
	}
}
