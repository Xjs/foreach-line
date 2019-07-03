package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func usage() {
	log.Println("Usage: foreach-line <list-of-things> <command to execute>")
	log.Println("Supply - to read list from stdin")
}

func run(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

const skipEmpty = true

func main() {
	printLines := false
	pattern := ""
	flag.BoolVar(&printLines, "print", printLines, "print each line before executing command")
	flag.StringVar(&pattern, "pattern", pattern, "replace this pattern by trimmed line content")
	flag.Parse()

	if len(flag.Args()) < 2 {
		usage()
		return
	}

	filename := flag.Arg(0)
	command := flag.Arg(1)
	cmdargs := flag.Args()[2:]

	var f io.Reader
	if filename == "-" {
		f = os.Stdin
	} else {
		var err error
		f, err = os.Open(filename)
		if err != nil {
			log.Fatalf("Couldn't open file %s for reading: %v\n", filename, err)
		}
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if skipEmpty && strings.TrimSpace(line) == "" {
			continue
		}
		if printLines {
			fmt.Println(line)
		}
		// make copy of arguments
		args := make([]string, len(cmdargs)+1)
		copy(args, cmdargs)
		if pattern != "" {
			for i := range args {
				args[i] = strings.ReplaceAll(args[i], pattern, line)
			}
		} else {
			args = append(args, line)
		}
		if err := run(command, args...); err != nil {
			log.Printf("Error running %s (args %s): %v\n", command, strings.Join(args, " "), err)
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v\n", err)
	}
}
