// Package types defines custom data structures used throughout the application.
// These types help in standardizing the data format, particularly for handling
// incoming messages in a more structured and convenient way.
package types

import (
	"time"

	"go.mau.fi/whatsmeow/types"
)

// Messages is a custom struct that represents a serialized WhatsApp message.
// It abstracts the raw message event from whatsmeow into a more accessible format,
// containing parsed information like the command, arguments, sender details, and more.
type Messages struct {
	// From is the JID (Jabber ID) of the chat where the message was sent.
	// This could be a group JID or a user JID.
	From types.JID
	// FromUser is the user part of the From JID.
	FromUser string
	// FromServer is the server part of the From JID (e.g., "s.whatsapp.net").
	FromServer string
	// FromMe indicates if the message was sent by the bot's own number.
	FromMe bool
	// ID is the unique identifier of the message.
	ID types.MessageID
	// IsGroup is true if the message was sent in a group chat.
	IsGroup bool
	// IsOwner is true if the message was sent by a user listed in the config.Owners.
	IsOwner bool
	// Sender is the JID of the user who sent the message. In a group, this is the
	// participant's JID; in a private chat, it's the same as From.
	Sender types.JID
	// SenderUser is the user part of the Sender JID.
	SenderUser string
	// SenderServer is the server part of the Sender JID.
	SenderServer string
	// Pushname is the display name of the sender as set in their WhatsApp profile.
	Pushname string
	// Timestamp is the time when the message was sent.
	Timestamp time.Time
	// Prefix is the command prefix used (e.g., '!' or '.').
	Prefix string
	// Command is the parsed command from the message text (e.g., "ping").
	Command string
	// Args is a slice of strings containing the arguments that followed the command.
	Args []string
	// Text is the complete, raw text content of the message.
	Text string
	// Body is an alias for Text, representing the main content of the message.
	Body string
}
