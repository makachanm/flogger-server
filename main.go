package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	messageChannel := make(chan Message, 100) // Buffer size 100
	terminateChannel := make(chan bool)

	server := NewFloggerServer(messageChannel, terminateChannel)
	server.StartServer()
	fmt.Println("Flogger server started...")

	// Goroutine to listen for messages from the server
	go func() {
		ui := NewLoggerUI(messageChannel)
		ui.Start()
	}()

	// Wait for interrupt signal to gracefully shut down the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down server...")
	terminateChannel <- true
	os.Remove("/tmp/flogger.sock")
	os.Exit(0)
}
