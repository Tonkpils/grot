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
}

func NewSlackAdapter(apiKey string) (*SlackAdapter, error) {
	api := slack.New(apiKey)
	_, err := api.AuthTest()
	if err != nil {
		return nil, err
	}
	rtm := api.NewRTM()
	return &SlackAdapter{
		RTM: rtm,
	}, err
}

// TODO: allow passing in an interface that defines options to send to the chat source
func (s *SlackAdapter) Send(res *grot.Response, msg string) error {
	msgParams := slack.NewPostMessageParameters()
	msgParams.AsUser = true
	msgParams.LinkNames = 1
	s.PostMessage(res.Room, msg, msgParams)
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
			case *slack.MessageEvent:
				// TODO: Msg could contain a Subtype such as message_changed
				// This means the original message was edited and the
				// SubMessage contains the new text
				bot.Logger.Printf("Message: %+v\n", ev)
				msg := grot.Message{
					User: grot.User{
						ID: ev.Msg.User,
					},
					Room: ev.Msg.Channel,
					Text: ev.Msg.Text,
				}
				bot.Receive(msg)
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
