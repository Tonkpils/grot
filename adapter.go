package grot

import (
	"log"
	"sync"
)

var (
	chatAdapterMu sync.RWMutex
	chatAdapter   ChatAdapter
)

func RegisterChatAdapter(adapter ChatAdapter) {
	chatAdapterMu.Lock()
	defer chatAdapterMu.Unlock()
	if chatAdapter != nil {
		log.Panic("chatbot: cannot load multiple chat adapters")
	}

	chatAdapter = adapter
}

type ChatAdapter interface {
	// Send is used to send a message to the chat source from the robot.
	// This will be used by listeners to respond to messages
	Send(message string) error
	// Receive ensures the adapter begins listening on the chat source.
	Receive(bot *Bot) error
}
