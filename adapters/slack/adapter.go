package slack

import (
	"log"
	"os"

	"github.com/Tonkpils/grot"
	"github.com/nlopes/slack"
)

func init() {
	apiKey := os.Getenv("GROT_SLACK_API_KEY")
	adapter, err := NewSlackAdapter(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	grot.RegisterChatAdapter(adapter)
}

type SlackAdapter struct {
	*slack.RTM
	// TODO: placing this here for convenience, ideally we'd pass
	// thse as optional arguments
	postParams slack.PostMessageParameters
}

func NewSlackAdapter(apiKey string) (*SlackAdapter, error) {
	api := slack.New(apiKey)
	_, err := api.AuthTest()
	if err != nil {
		return nil, err
	}
	rtm := api.NewRTM()
	msgParams := slack.NewPostMessageParameters()
	msgParams.AsUser = true
	msgParams.LinkNames = 1
	return &SlackAdapter{
		RTM:        rtm,
		postParams: msgParams,
	}, err
}

// TODO: allow passing in an interface that defines options to send to the chat source
func (s *SlackAdapter) Send(msg string) error {
	s.PostMessage("#bender-playground", msg, s.postParams)
	return nil
}

func (s *SlackAdapter) Receive(bot *grot.Bot) error {
	go s.ManageConnection()

Loop:
	for {
		select {
		case msg := <-s.IncomingEvents:
			bot.Logger.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// TODO: here we could tell bot that we are connected
			case *slack.ConnectedEvent:
				bot.Logger.Printf("Infos: %+v\n", ev.Info)
				s.Send("Bite my shiny metal ass!")
			case *slack.MessageEvent:
				bot.Logger.Printf("Message: %+v\n", ev)
				// TODO: here we would send the bot the message so listeners can handle them
			case *slack.PresenceChangeEvent:
				bot.Logger.Printf("Presence Change: %+v\n", ev)
			case *slack.RTMError:
				bot.Logger.Printf("Error: %s\n", ev.Error())
			case *slack.InvalidAuthEvent:
				bot.Logger.Printf("Invalid credentials")
				break Loop
			default:
				bot.Logger.Printf("Unexpected: %+v\n", msg.Data)
			}
		}
	}

	return nil
}
