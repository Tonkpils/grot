package grot

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

type Bot struct {
	Client ChatAdapter
	Logger *log.Logger
	Router *http.ServeMux
}

func NewBot() *Bot {
	logger := log.New(os.Stdout, "", 0)
	router := http.NewServeMux()

	return &Bot{
		Client: chatAdapter,
		Logger: logger,
		Router: router,
	}
}

func (b *Bot) Run() error {
	go b.Client.Receive(b)

	b.Router.HandleFunc("/grot/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PONG"))
	})

	// TODO: create routing structure to allow providing handlers
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), b.Router); err != nil {
		return err
	}

	return nil
}
