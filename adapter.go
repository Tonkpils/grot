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
	// Takes a Response which contains metadata about the message origin
	// and a message string.
	Send(*Response, string) error
	// Receive ensures the adapter begins listening on the chat source.
	Receive(bot *Bot) error
}
