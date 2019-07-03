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

const skipEmpty = true

func readerFrom(name string) (io.Reader, error) {
	if name == "-" {
		return os.Stdin, nil
	}
	return os.Open(name)
}

func main() {
	printLines := false
	skipIOFail := false
	pattern := ""
	stdin := ""
	stdout := ""
	stderr := ""
	flag.BoolVar(&printLines, "print", printLines, "print each line before executing command")
	flag.StringVar(&pattern, "pattern", pattern, "replace this pattern by trimmed line content")
	flag.StringVar(&stdin, "stdin", stdin, "read stdin for command from this file (pattern replacement will be done). Supply - for this command's stdin (obviously incompatible with reading lines from stdin)")
	flag.StringVar(&stdout, "stdout", stdout, "redirect stdout for command from this file (pattern replacement will be done)")
	flag.StringVar(&stderr, "stderr", stderr, "redirect stderr for command from this file (pattern replacement will be done)")
	flag.BoolVar(&skipIOFail, "skip-io-fail", skipIOFail, "skip line's command execution if attaching stdin/stdout/stderr files fails (if false, continue with this command's stdout/stderr and no stdin)")
	flag.Parse()

	if len(flag.Args()) < 2 {
		usage()
		return
	}

	filename := flag.Arg(0)
	command := flag.Arg(1)
	cmdargs := flag.Args()[2:]

	f, err := readerFrom(filename)
	if err != nil {
		log.Fatalf("Couldn't open file %s for reading: %v\n", filename, err)
	}

	lineno := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lineno++
		if skipEmpty && strings.TrimSpace(line) == "" {
			continue
		}
		if printLines {
			fmt.Println(line)
		}
		// make copy of arguments
		args := make([]string, len(cmdargs), len(cmdargs)+1)
		copy(args, cmdargs)
		stdin := stdin
		stdout := stdout
		stderr := stderr

		if pattern != "" {
			for i := range args {
				args[i] = strings.ReplaceAll(args[i], pattern, line)
			}
			stdin = strings.ReplaceAll(stdin, pattern, line)
			stdout = strings.ReplaceAll(stdout, pattern, line)
			stderr = strings.ReplaceAll(stderr, pattern, line)
		} else {
			args = append(args, line)
		}

		c := exec.Command(command, args...)

		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if stdin == "-" && filename != "-" {
			c.Stdin = os.Stdin
		} else if stdin != "" {
			in, err := readerFrom(stdin)
			if err != nil {
				log.Printf("Error opening %s for reading (line %d): %v\n", stdin, lineno, err)
				if skipIOFail {
					continue
				}
			} else {
				c.Stdin = in
			}
		}

		if stdout != "" {
			out, err := os.Create(stdout)
			if err != nil {
				log.Printf("Error opening %s for writing as stdout (line %d): %v\n", stdout, lineno, err)
				if skipIOFail {
					continue
				}
			}
			c.Stdout = out
		}

		if stderr != "" {
			out, err := os.Create(stderr)
			if err != nil {
				log.Printf("Error opening %s for writing as stderr (line %d): %v\n", stderr, lineno, err)
				if skipIOFail {
					continue
				}
			}
			c.Stderr = out
		}

		if err := c.Run(); err != nil {
			log.Printf("Error running %s (args %s): %v\n", command, strings.Join(args, " "), err)
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v\n", err)
	}
}
