package main

import (
	"./commands"
	"./engine"
	"bufio"
	"log"
	"os"
)

func main() {
	inputFile := "file.txt"
	eventLoop := new(engine.Loop)
	eventLoop.Start()
	if input, err := os.Open(inputFile); err == nil {
		defer input.Close()
		scanner := bufio.NewScanner(input)
		for scanner.Scan() {
			commandLine := scanner.Text()
			cmd := commands.Parse(commandLine)
			eventLoop.Post(cmd)
		}
	} else {
		log.Fatal("Can't open a file" + err.Error())
	}
	eventLoop.AwaitFinish()
}
