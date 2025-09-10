// Package main is the entry point for the WhatsApp bot application.
// It initializes the bot client, connects to WhatsApp, and handles graceful shutdown.
package main

import (
	"botwa/client"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// main is the primary function that starts the application.
// It sets up the WhatsApp client and listens for system signals (like Ctrl+C)
// to disconnect the client gracefully.
func main() {
	// Initialize the WhatsApp client, which sets up the database connection,
	// logs in, and registers the event handler.
	client.Init()

	// Create a channel to listen for termination signals.
	// This allows the application to shut down cleanly when it receives
	// an interrupt (os.Interrupt) or a termination signal (syscall.SIGTERM).
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block execution until a signal is received on the 'stop' channel.
	<-stop

	// Once a signal is received, log a shutdown message and disconnect the client.
	fmt.Println("Shutting down the bot.")
	client.WhatsAppClient.Disconnect()
}
