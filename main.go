package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	log.SetFlags(0)
	var eventsPath string
	flag.StringVar(&eventsPath, "events", "events.yaml", "YAML file defining events")
	var token string
	flag.StringVar(&token, "token", "", "Telegram bot token")
	var chatID int64
	flag.Int64Var(&chatID, "chat-id", 0, "Telegram chat ID")
	flag.Parse()

	events, err := readEvents(eventsPath)
	if err != nil {
		log.Fatalln(err)
	}

	bot, err := newBot(token, chatID)
	if err != nil {
		log.Fatalln(err)
	}
	go bot.listen()

	for t := range time.Tick(time.Hour) {
		if t.Round(time.Hour).Hour() != 9 { // run once a day between 8:30 and 9:29
			continue
		}
		bot.send(events.today()...)
	}
}
