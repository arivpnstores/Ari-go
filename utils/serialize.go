// Package utils provides helper functions and utilities for the bot.
// This file, serialize.go, is responsible for converting raw message events
// from the whatsmeow library into a structured, more usable custom format.
package utils

import (
	"botwa/config"
	"botwa/types"
	"strings"
	"time"

	"go.mau.fi/whatsmeow/types/events"
)

// Serialize converts a raw *events.Message into a custom types.Messages struct.
// This process involves extracting key information, parsing for commands and arguments,
// and standardizing the data for easier use throughout the application.
//
// Parameters:
//   ctx: A pointer to the raw message event (*events.Message) from whatsmeow.
//
// Returns:
//   A types.Messages struct populated with structured information from the raw event.
func Serialize(ctx *events.Message) types.Messages {
	// Set the timezone to Asia/Jakarta for consistent timestamps.
	loc, _ := time.LoadLocation("Asia/Jakarta")

	// Extract the message body and check for a command prefix.
	body := GetText(ctx)
	prefix := GetPrefix(body)
	words := strings.Fields(body)
	info := ctx.Info

	// Parse the command and its arguments from the message body.
	cmd, args := "", []string{}
	if len(words) > 0 && strings.HasPrefix(words[0], prefix) {
		cmd = strings.TrimPrefix(words[0], prefix)
		args = words[1:]
	}

	// Construct and return the serialized message struct.
	return types.Messages{
		From:         info.Chat,
		FromUser:     info.Chat.User,
		FromServer:   info.Chat.Server,
		FromMe:       info.IsFromMe,
		ID:           info.ID,
		IsGroup:      info.IsGroup,
		IsOwner:      isOwner(info.Sender.User), // Check if the sender is an owner.
		Sender:       info.Sender,
		SenderUser:   info.Sender.User,
		SenderServer: info.Sender.Server,
		Pushname:     info.PushName,
		Timestamp:    info.Timestamp.In(loc),
		Prefix:       prefix,
		Command:      cmd,
		Args:         args,
		Text:         strings.Join(args, " "), // The text part of the message (excluding the command).
		Body:         body,                    // The full, original message text.
	}
}

// isOwner checks if a given user JID belongs to a bot owner.
// It compares the user's ID against the list of owner IDs in the config.
//
// Parameters:
//   user: The user part of a JID (e.g., "6281234567890").
//
// Returns:
//   A boolean value: true if the user is an owner, false otherwise.
func isOwner(user string) bool {
	for _, owner := range config.Owners {
		if owner == user {
			return true
		}
	}
	return false
}
