// Package utils provides helper functions and utilities for the bot.
// This file, extract.go, contains functions for extracting specific pieces
// of information from raw message events, such as the message text or command prefix.
package utils

import (
	"botwa/config"

	"go.mau.fi/whatsmeow/types/events"
)

// GetText extracts the primary text content from a message event.
// It intelligently checks different message types (e.g., image/video captions,
// extended text, or a simple conversation) and returns the relevant text.
//
// Parameters:
//   ctx: A pointer to the raw message event (*events.Message) from whatsmeow.
//
// Returns:
//   A string containing the extracted text content. Returns an empty string if
//   no text can be found.
func GetText(ctx *events.Message) string {
	if ctx == nil || ctx.Message == nil {
		return ""
	}

	msg := ctx.Message

	switch {
	case msg.ImageMessage != nil && msg.ImageMessage.Caption != nil:
		return *msg.ImageMessage.Caption
	case msg.VideoMessage != nil && msg.VideoMessage.Caption != nil:
		return *msg.VideoMessage.Caption
	case msg.ExtendedTextMessage != nil && msg.ExtendedTextMessage.Text != nil:
		return *msg.ExtendedTextMessage.Text
	case msg.DocumentMessage != nil && msg.DocumentMessage.Caption != nil:
		return *msg.DocumentMessage.Caption
	case msg.Conversation != nil:
		return *msg.Conversation
	default:
		return ""
	}
}

// GetPrefix checks if a given text starts with one of the recognized command prefixes.
// The list of valid prefixes is defined in the config package.
//
// Parameters:
//   text: The string to check for a prefix.
//
// Returns:
//   A string containing the detected prefix if found. Returns an empty string if
//   the text does not start with a valid prefix.
func GetPrefix(text string) string {
	if len(text) == 0 {
		return ""
	}

	first := rune(text[0])
	for _, prefix := range config.Prefixes {
		if first == prefix {
			return string(first)
		}
	}

	return ""
}
