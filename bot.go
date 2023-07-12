package main

import (
	"fmt"
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const dataDone = "done"

type bot struct {
	api    *tg.BotAPI
	chatID int64
}

func newBot(token string, chatID int64) (*bot, error) {
	api, err := tg.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	log.Println("INFO", "Authorized as", api.Self.UserName)
	return &bot{api: api, chatID: chatID}, nil
}

func (b *bot) send(e event) error {
	msg := tg.NewMessage(b.chatID, e.String())
	msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Done", dataDone),
		),
	)
	_, err := b.api.Send(msg)
	return fmt.Errorf("failed to send message: %w", err)
}

func (b *bot) listen() {
	u := tg.NewUpdate(0)
	u.Timeout = 60
	for update := range b.api.GetUpdatesChan(u) {
		if update.FromChat().ID != b.chatID {
			continue
		}
		switch {
		case update.CallbackQuery != nil:
			if err := b.handleCallback(update); err != nil {
				log.Println("ERROR", err)
				continue
			}
		case update.Message.IsCommand():
			switch update.Message.Command() {
			case "next":
				// TODO:
			}
		}
	}
}

func (b *bot) handleCallback(update tg.Update) error {
	if _, err := b.api.Request(tg.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)); err != nil {
		return fmt.Errorf("failed to receive callback: %w", err)
	}
	switch update.CallbackData() {
	case dataDone:
		return b.callbackDone(update.CallbackQuery)
	}
	return nil
}

func (b *bot) callbackDone(cq *tg.CallbackQuery) error {
	edit := tg.NewEditMessageText(b.chatID, cq.Message.MessageID, fmt.Sprintf(formatDone, cq.From.UserName, cq.Message.Text))
	edit.ParseMode = tg.ModeMarkdown
	_, err := b.api.Send(edit)
	return fmt.Errorf("failed to edit message: %w", err)
}
