// Package client handles the initialization, connection, and management of the WhatsApp client.
// It is responsible for setting up the database, handling device sessions, and managing the
// connection lifecycle, including QR code generation for new logins.
package client

import (
	"botwa/handler"
	"context"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3" // SQLite3 driver for the database
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// WhatsAppClient is a globally accessible instance of the whatsmeow client.
// It is initialized by the Init function and used throughout the application
// to interact with the WhatsApp API.
var WhatsAppClient *whatsmeow.Client

// Init initializes the primary WhatsApp client.
// This function performs the following steps:
// 1. Sets up a logger for the client.
// 2. Creates a new SQL-based store (using SQLite3) to manage session data.
// 3. Fetches the first available device from the store or creates a new one.
// 4. Initializes the whatsmeow client with the device and logger.
// 5. Registers the main event handler to process incoming events.
// 6. Connects to WhatsApp. If it's the first time, it generates a QR code
//    for login. Otherwise, it attempts to reconnect using the saved session.
func Init() {
	log := waLog.Stdout("Client", "INFO", true)

	// Create a container for the session data using SQLite3.
	// The session data is stored in "session.db".
	container, err := sqlstore.New(context.Background(), "sqlite3", "file:session.db?_foreign_keys=on", log)
	if err != nil {
		log.Errorf("DB error: %v", err)
		return
	}

	// Get the first device from the store. If no device is found, a new one will be created.
	device, err := container.GetFirstDevice(context.Background())
	if err != nil {
		log.Errorf("Device error: %v", err)
		return
	}

	// Initialize the whatsmeow client with the retrieved device and logger.
	WhatsAppClient = whatsmeow.NewClient(device, log)
	
	// Register the global event handler from the handler package.
	WhatsAppClient.AddEventHandler(func(evt interface{}) {
		handler.EventHandler(evt, WhatsAppClient)
	})

	// Check if the client is already logged in by looking for a stored ID.
	if WhatsAppClient.Store.ID == nil {
		// If not logged in, get the QR channel and connect.
		qrChan, _ := WhatsAppClient.GetQRChannel(context.Background())
		err := WhatsAppClient.Connect()
		if err != nil {
			log.Errorf("Connect error: %v", err)
			return
		}

		// Listen on the QR channel for events.
		for evt := range qrChan {
			switch evt.Event {
			case "code":
				// A new QR code is available. Save it to a file.
				qrcode.WriteFile(evt.Code, qrcode.Medium, 256, "qrcode.png")
				fmt.Println("QR code saved to qrcode.png")
			case "success", "timeout":
				// Login was successful or the QR code timed out. Remove the QR code file.
				os.Remove("qrcode.png")
			}
		}
	} else {
		// If already logged in, just connect.
		err := WhatsAppClient.Connect()
		if err != nil {
			log.Errorf("Connect error: %v", err)
			return
		}
		// Clean up any old QR code file that might exist.
		os.Remove("qrcode.png")
	}
}
