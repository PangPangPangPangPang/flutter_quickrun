// Package main provides ...
package main

import (
	// "encoding/json"
	"bufio"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func main() {
	flutterCmd := "run"
	for idx, args := range os.Args {
		switch idx {
		case 1:
			if args == "run" || args == "attach" {
				flutterCmd = args
			} else {
				log.Fatal("error args")
			}
		}
	}
	defaultUrl := "http://127.0.0.1:9100/"
	myCmd := exec.Command("bash", "-c", "pub global run devtools")
	execCommand(myCmd, func(msg string) {
		r, _ := regexp.Match("DevTools", []byte(msg))
		if r {
			words := strings.Fields(msg)
			for idx := range words {
				r, _ := regexp.Match("http://", []byte(words[idx]))
				if r {
					defaultUrl = words[idx]
				}
			}
		}
	})
	runCmd := exec.Command("bash", "-c", fmt.Sprintln("flutter "+flutterCmd))
	stdIn, _ := runCmd.StdinPipe()
	execCommand(runCmd, func(msg string) {
		fmt.Println(msg)
		r, _ := regexp.Match("is available at", []byte(msg))
		if r {
			words := strings.Fields(msg)
			for e := range words {
				r, _ := regexp.Match("http://", []byte(words[e]))
				if r {
					open := fmt.Sprintf("%s %s?uri=%s", openCmd(), defaultUrl, url.QueryEscape(words[e]))
					runCmd := exec.Command("bash", "-c", open)
					execCommand(runCmd, nil)
				}
			}
		}
	})
	var input string
	for true {
		fmt.Scanln(&input)
		io.WriteString(stdIn, input)
	}

}

type Handler func(msg string)

func execCommand(c *exec.Cmd, handler Handler) {
	readPipe := func(reader io.Reader, prefix string) {
		r := bufio.NewReader(reader)
		var outStr string
		var line []byte
		for true {
			line, _, _ = r.ReadLine()
			if line != nil {
				outStr = string(line)
				if handler != nil {
					handler(outStr)
				}
			}
		}
	}
	cmdOut, _ := c.StdoutPipe()
	cmdErr, _ := c.StderrPipe()
	go readPipe(cmdOut, "Output: ")
	go readPipe(cmdErr, "Error: ")
	c.Start()

}

func openCmd() string {
	if runtime.GOOS == "darwin" {
		return "open"
	} else {
		return "xdg-open"
	}
}
