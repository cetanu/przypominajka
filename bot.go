package main

import (
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type bot struct {
	api    *tg.BotAPI
	chatID int64
}

const donePayload = "done"

func newBot(token string, chatID int64) (*bot, error) {
	api, err := tg.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	log.Println("INFO", "Authorized as", api.Self.UserName)

	return &bot{api, chatID}, nil
}

func (b *bot) notify(message string) error {
	msg := tg.NewMessage(b.chatID, message)
	msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Done", donePayload),
		),
	)
	_, err := b.api.Send(msg)
	return err
}

func (b *bot) edit(messageID int, message string) error {
	_, err := b.api.Send(tg.NewEditMessageText(b.chatID, messageID, message))
	return err
}

func (b *bot) listen() {
	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			callback := tg.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := b.api.Request(callback); err != nil {
				log.Println("ERROR", err)
				continue
			}

			if update.CallbackQuery.Data != donePayload {
				continue
			}

			if err := b.edit(update.CallbackQuery.Message.MessageID, "po robocie"); err != nil {
				log.Println("ERROR", err)
				continue
			}
		}
	}
}
