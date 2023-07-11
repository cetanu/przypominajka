package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	log.SetFlags(0)
	var token string
	flag.StringVar(&token, "token", "", "Telegram bot token")
	var chatID int64
	flag.Int64Var(&chatID, "chat-id", 87974246, "Telegram chat ID")
	flag.Parse()

	bot, err := newBot(token, chatID)
	if err != nil {
		log.Fatalln(err)
	}
	go bot.listen()

	bdays, err := readBirthdays()
	if err != nil {
		log.Fatalln(err)
	}

	for t := range time.Tick(time.Hour) {
		if t.Round(time.Hour).Hour() != 9 { // run once a day between 8:30 and 9:29
			continue
		}
		bot.notify(bdays.at(t)...)
	}
}
