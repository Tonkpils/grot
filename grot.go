package grot

import (
	"log"
	"net/http"
	"os"
)

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
