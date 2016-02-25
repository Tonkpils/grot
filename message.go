package grot

type User struct {
	ID   string
	Name string
}

type Message struct {
	User User
	Room string
	Text string
}
