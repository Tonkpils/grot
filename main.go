package main

import (
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
}

func main() {
	log.Println("system going online...")
	bot := NewBot()

	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}
