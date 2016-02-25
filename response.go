package grot

type Response struct {
	Message
	Matches []string
	bot     *Bot
}

// Send a message back to the room where the response originated from.
func (r *Response) Send(msg string) {
	r.bot.Client.Send(r, msg)
}
