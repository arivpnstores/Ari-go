// Package handler provides functions for processing events received from the WhatsApp client.
// It acts as the central point for routing different event types to their respective handlers.
package handler

import (
	"botwa/utils"
	"fmt"

	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// EventHandler is the primary event handler for the WhatsApp client.
// It receives all events from the whatsmeow client, determines their type,
// and delegates them to the appropriate logic.
//
// Parameters:
//   evt: The event interface{} received from the client. This can be any event type
//        defined by the whatsmeow library (e.g., *events.Message, *events.Receipt, etc.).
//   client: A pointer to the active whatsmeow.Client instance, used to perform actions
//           like sending messages or marking them as read.
func EventHandler(evt interface{}, client *whatsmeow.Client) {
	switch v := evt.(type) {
	// Case for handling incoming messages.
	case *events.Message:
		// Serialize the raw message event into a more manageable custom format.
		m := utils.Serialize(v)

		// Ignore messages originating from newsletters to prevent processing channel updates.
		if m.FromServer == "newsletter" {
			return
		}

		// Automatically mark status updates as read.
		if m.From.String() == "status@broadcast" {
			err := client.MarkRead(
				[]waTypes.MessageID{m.ID},
				m.Timestamp,
				m.From,
				m.Sender,
			)
			if err != nil {
				// Log if marking the status as read fails.
				fmt.Println("Failed to mark status as read:", err)
			}
		}

		// Pass the serialized message to the command handler for further processing.
		utils.HandleCommand(client, m, v)
	}
}
