package main

import (
	"fmt"
	"log"
	"time"

	"git.sr.ht/~tymek/przypominajka/format"
	"git.sr.ht/~tymek/przypominajka/models"
	"git.sr.ht/~tymek/przypominajka/storage"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const dataDone = "done"

type bot struct {
	api    *tg.BotAPI
	chatID int64
	s      storage.Interface
}

func newBot(token string, chatID int64, s storage.Interface) (*bot, error) {
	api, err := tg.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	log.Println("INFO", "Authorized as", api.Self.UserName)
	return &bot{api: api, chatID: chatID, s: s}, nil
}

func (b *bot) send(e models.Event) error {
	msg := tg.NewMessage(b.chatID, e.Format(false))
	msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Done", dataDone),
		),
	)
	_, err := b.api.Send(msg)
	return fmt.Errorf("failed to send message: %w", err)
}

func (b *bot) serve() {
	for t := range time.Tick(time.Hour) {
		if t.Round(time.Hour).Hour() != 9 { // run once a day between 8:30 and 9:29
			continue
		}
		events, err := storage.Today(b.s)
		if err != nil {
			log.Println("ERROR", err)
			continue
		}
		for _, e := range events {
			if err := b.send(e); err != nil {
				log.Println("ERROR", err)
			}
		}
	}
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
			if err := b.handleCommand(update); err != nil {
				log.Println("ERROR", err)
				continue
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
	edit := tg.NewEditMessageText(b.chatID, cq.Message.MessageID, fmt.Sprintf(format.MsgDone, cq.From.UserName, cq.Message.Text))
	edit.ParseMode = tg.ModeMarkdown
	_, err := b.api.Send(edit)
	return fmt.Errorf("failed to edit message: %w", err)
}

func (b *bot) handleCommand(update tg.Update) error {
	// NOTE: if another bot has /next command, then this will be triggered.
	// To prevent this behavior, we can CommandWithAt() and check whether
	// <command>@<bot_name> matches.
	switch update.Message.Command() {
	case "next":
		text := format.MsgNoEvents
		events, err := storage.Next(b.s)
		if err != nil {
			return err
		}
		if len(events) > 0 {
			text = events.Format(events[0].Month, events[0].Day)
		}
		if _, err := b.api.Send(tg.NewMessage(b.chatID, text)); err != nil {
			return err
		}
	}
	return nil
}
