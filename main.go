package main

import (
	"bufio"
	"flag"
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
	flag.BoolVar(&printLines, "print", printLines, "print each line before executing command")
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
		if printLine {
			fmt.Println(line)
		}
		input := append(cmdargs, line)
		if err := run(command, input...); err != nil {
			log.Printf("Error running %s (args %s): %v\n", command, strings.Join(input, " "), err)
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v\n", err)
	}
}
