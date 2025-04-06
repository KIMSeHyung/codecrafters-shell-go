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

func qoutesProcess(str string) []string {
	var sb strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	previous := rune(' ')
	var cache []string

	for _, r := range str {
		switch {
		case r == '\'' && !inDoubleQuote && !inSingleQuote:
			inSingleQuote = true

		case r == '"' && !inDoubleQuote && !inSingleQuote:
			inDoubleQuote = true

		case r == '\'' && inSingleQuote:
			cache = append(cache, sb.String())
			sb.Reset()
			inSingleQuote = false

		case r == '"' && inDoubleQuote:
			cache = append(cache, sb.String())
			sb.Reset()
			inDoubleQuote = false

		case r == ' ' && previous != ' ' && !inSingleQuote && !inDoubleQuote:
			sb.WriteRune(' ')

		case inSingleQuote || inDoubleQuote || r != ' ':
			sb.WriteRune(r)
		}

		previous = r
	}

	if sb.Len() > 0 {
		cache = append(cache, sb.String())
	}
	return cache
}

func main() {
	builtinCommands := map[string]bool{
		"exit": true, "echo": true, "type": true, "pwd": true, "cd": true,
	}
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
				str := strings.Join(cmd[1:], " ")
				cache := qoutesProcess(str)
				fmt.Println(strings.Join(cache, ""))

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

			case "cd":
				inputPath := cmd[1]
				if inputPath == "~" {
					home, _ := os.UserHomeDir()
					os.Chdir(home)
					break
				}
				_, err := os.Stat(inputPath)
				if err != nil {
					fmt.Printf("cd: %s: No such file or directory\n", inputPath)
				} else {
					os.Chdir(inputPath)
				}

			case "cat":
				str := strings.Join(cmd[1:], " ")
				// var sb strings.Builder
				// isCompleteQuote := false
				// prevQuote := false
				// var tmp []string
				// for _, x := range str {
				// 	if x == '\'' {
				// 		isCompleteQuote = prevQuote && !isCompleteQuote
				// 		prevQuote = !prevQuote
				// 	}
				// 	sb.WriteRune(x)
				// 	if isCompleteQuote {
				// 		s := strings.ReplaceAll(sb.String(), "'", "")
				// 		s = strings.ReplaceAll(s, " ", "")
				// 		tmp = append(tmp, strings.TrimSpace(s))
				// 		isCompleteQuote = !isCompleteQuote
				// 		sb.Reset()
				// 	}
				// }
				// var catCmd []string
				// if len(tmp) > 0 {
				// 	catCmd = tmp
				// } else {
				// 	catCmd = append(catCmd, strings.Join(strings.Fields(str), " "))
				// }
				cache := qoutesProcess(str)
				var catCmd []string
				for _, x := range cache {
					catCmd = append(catCmd, strings.TrimSpace(x))
				}

				c := exec.Command(cmd[0], catCmd...)
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				c.Run()

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
