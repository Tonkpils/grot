package grot

import (
	"log"
	"net/http"
	"os"
	"regexp"
)

type Bot struct {
	Logger    *log.Logger
	Router    *http.ServeMux
	Client    ChatAdapter
	listeners []Listener
}

func NewBot() *Bot {
	logger := log.New(os.Stdout, "", 0)
	router := http.NewServeMux()

	// TODO: if chatAdapter is not set, simply allow http routing

	return &Bot{
		Client: chatAdapter,
		Logger: logger,
		Router: router,
	}
}

func (b *Bot) Hear(pattern string, fn func(res *Response)) {
	l := &listener{
		fn:    fn,
		regex: regexp.MustCompile(pattern),
	}
	b.listeners = append(b.listeners, l)
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

func (b *Bot) Receive(msg Message) error {
	b.Logger.Printf("Robot received message %+v\n", msg)
	for _, l := range b.listeners {
		res := &Response{
			bot:     b,
			Message: msg,
		}
		l.Handle(res)
	}

	return nil
}
