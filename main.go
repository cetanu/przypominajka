package main

import (
	"flag"
	"log"
	"sync"
)

func main() {
	log.SetFlags(0)
	var token string
	flag.StringVar(&token, "token", "", "Telegram bot token")
	var chatID int64
	flag.Int64Var(&chatID, "chat-id", 87974246, "Telegram chat ID")
	flag.Parse()

	bdays, err := readBirthdays()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(bdays)

	bot, err := newBot(token, chatID)
	if err != nil {
		log.Fatalln(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go bot.listen()
	log.Println(bot.notify("Hello world!"))
	wg.Wait()
}
