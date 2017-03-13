package main

import (
	"os"
	"log"
)

var apiToken = os.Getenv("SLACK_API_TOKEN")
var bot *SlackBot
var id int = 0

func main() {
	bot = &SlackBot{ApiToken: apiToken, EventHandler: handler}
	bot.Run()
}


func handler(event *Event) {
	log.Printf("Incoming message: %+v\n", event)
	//bytes, _ := json.Marshal(event)
	if event.Type == "message" {
		response := &Event{
			Id: id,
			Type: "message",
			Channel:event.Channel,
			Text:"Answer from bot",
		}
		id++
		log.Printf("Outcoming message: %+v\n", response)
		bot.Send(response)
	}
}
