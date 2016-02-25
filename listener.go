package grot

import "regexp"

type Listener interface {
	Handle(*Response)
}

type listener struct {
	regex *regexp.Regexp
	fn    func(msg *Response)
}

func (l *listener) Handle(res *Response) {
	if !l.regex.MatchString(res.Text) {
		return
	}

	res.Matches = l.regex.FindAllStringSubmatch(res.Text, -1)[0]
	l.fn(res)
}
