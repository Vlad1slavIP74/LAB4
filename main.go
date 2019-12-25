package main

import (
	"./commands"
	"./engine"
	"bufio"
	"log"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	inputFile := "file.txt"
	eventLoop := new(engine.Loop)
	eventLoop.Start(&wg)
	if input, err := os.Open(inputFile); err == nil {
		defer input.Close()
		scanner := bufio.NewScanner(input)
		for scanner.Scan() {
			wg.Add(2)
			commandLine := scanner.Text()
			cmd := commands.Parse(commandLine)
			eventLoop.Post(cmd, &wg)
		}
	} else {
		log.Fatal("Can't open a file" + err.Error())
	}
	wg.Wait()
}
