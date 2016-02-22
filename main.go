package main

import (
	"log"
	"net/http"
	"os"
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
		log.Panic("grot: cannot load multiple chat adapters")
	}

	chatAdapter = adapter
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
}

type ChatAdapter interface {
	// Send is used to send a message to the chat source from the robot.
	// This will be used by listeners to respond to messages
	Send(message string) error
	// Receive ensures the adapter begins listening on the chat source.
	Receive(bot *Bot) error
}

type Bot struct {
	Client ChatAdapter
	Logger *log.Logger
}

func NewBot() *Bot {
	logger := log.New(os.Stdout, "", 0)
	return &Bot{
		Client: chatAdapter,
		Logger: logger,
	}
}

func (b *Bot) Run() error {
	go b.Client.Receive(b)

	// TODO: create routing structure to allow providing handlers
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		return err
	}

	return nil
}

func main() {
	log.Println("system going online...")
	// Configuration ....
	// - Chat adapter
	// - Additional command adapters
	bot := NewBot()

	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}
