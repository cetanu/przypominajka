package main

import (
	"fmt"
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
	return &bot{api: api, chatID: chatID}, nil
}

func (b *bot) notify(infos ...info) {
	for _, i := range infos {
		msg := tg.NewMessage(b.chatID, i.String())
		msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Done", donePayload),
			),
		)
		if _, err := b.api.Send(msg); err != nil {
			log.Println("ERROR", err)
		}
	}
}

func (b *bot) listen() {
	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			if _, err := b.api.Request(tg.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)); err != nil {
				log.Println("ERROR", err)
				continue
			}

			if update.CallbackQuery.Data != donePayload {
				continue
			}

			edit := tg.NewEditMessageText(
				b.chatID,
				update.CallbackQuery.Message.MessageID,
				fmt.Sprintf(
					"%s\n\n_%s złożył(a) życzenia_",
					update.CallbackQuery.Message.Text,
					update.CallbackQuery.From.UserName,
				),
			)
			edit.ParseMode = tg.ModeMarkdown
			if _, err := b.api.Send(edit); err != nil {
				log.Println("ERROR", err)
				continue
			}
		}
	}
}
