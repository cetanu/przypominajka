package bot

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"git.sr.ht/~tymek/przypominajka/bot/wizard"
	"git.sr.ht/~tymek/przypominajka/format"
	"git.sr.ht/~tymek/przypominajka/models"
	"git.sr.ht/~tymek/przypominajka/storage"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const dataNotifyDone = "done"

type Bot struct {
	mu      sync.RWMutex
	api     *tg.BotAPI
	chatID  int64
	s       storage.Interface
	wizards map[string]wizard.Interface
	consume wizard.Consume
}

func New(token string, chatID int64, s storage.Interface, wizards ...wizard.Interface) (*Bot, error) {
	api, err := tg.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	log.Println("INFO", "Authorized as", api.Self.UserName)
	m := make(map[string]wizard.Interface, len(wizards))
	for _, w := range wizards {
		m[w.Name()] = w
	}
	return &Bot{api: api, chatID: chatID, s: s, wizards: m}, nil
}

func ListenAndServe(token string, chatID int64, s storage.Interface) error {
	b, err := New(token, chatID, s, &wizard.Add{})
	if err != nil {
		return err
	}
	go b.Listen()
	b.Serve()
	return nil
}

func (b *Bot) Notify(e models.Event) error {
	msg := tg.NewMessage(b.chatID, e.Format(false))
	msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Done", dataNotifyDone),
		),
	)
	return b.send(msg)
}

func (b *Bot) Listen() {
	u := tg.NewUpdate(0)
	u.Timeout = 60
	for update := range b.api.GetUpdatesChan(u) {
		go func(update tg.Update) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("PANIC", r)
					debug.PrintStack()
				}
			}()
			if err := b.handle(update); err != nil {
				log.Println("ERROR", err)
			}
		}(update)
	}
}

func (b *Bot) Serve() {
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
			if err := b.Notify(e); err != nil {
				log.Println("ERROR", err)
			}
		}
	}
}

func (b *Bot) handle(update tg.Update) error {
	if chat := update.FromChat(); chat == nil || chat.ID != b.chatID {
		return nil
	}
	switch {
	case update.CallbackQuery != nil:
		if _, err := b.api.Request(tg.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)); err != nil {
			return err
		}
		switch data := update.CallbackQuery.Data; data {
		case dataNotifyDone:
			return b.handleCallbackNotifyDone(update.CallbackQuery)
		default:
			name, _, ok := strings.Cut(update.CallbackData(), wizard.CallbackSep)
			if !ok {
				return nil
			}
			w, ok := b.wizards[name]
			if !ok {
				return nil
			}
			return b.runConsume(w.Next, update)
		}

	case update.Message.IsCommand():
		// NOTE: if another bot has /next command, then this will be triggered.
		// To prevent this behavior, we can CommandWithAt() and check whether
		// <command>@<bot_name> matches.
		switch cmd := update.Message.Command(); cmd {
		case "abort":
			b.mu.Lock()
			for _, w := range b.wizards {
				w.Reset()
			}
			b.mu.Unlock()
		case "next":
			return b.handleCommandNext(update)
		default:
			if w, ok := b.wizards[update.Message.Command()]; ok {
				return b.send(w.Start(update))
			}
		}

	case update.Message.Text != "":
		return b.runConsume(b.consume, update)
	}
	return nil
}

func (b *Bot) handleCallbackNotifyDone(cq *tg.CallbackQuery) error {
	edit := tg.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, fmt.Sprintf(format.MessageDone, cq.From.UserName, cq.Message.Text))
	edit.ParseMode = tg.ModeMarkdown
	_, err := b.api.Send(edit)
	return fmt.Errorf("failed to edit message: %w", err)
}

func (b *Bot) handleCommandNext(update tg.Update) error {
	events, err := storage.Next(b.s)
	if err != nil {
		return err
	}
	return b.send(tg.NewMessage(update.FromChat().ID, events.String()))
}

func (b *Bot) send(c tg.Chattable) error {
	if c == nil {
		log.Println("WARN", "nil message")
		return nil
	}
	if _, err := b.api.Send(c); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func (b *Bot) runConsume(c wizard.Consume, update tg.Update) error {
	if c == nil {
		return nil
	}
	b.mu.RLock()
	msg, consume, err := c(b.s, update)
	b.mu.RUnlock()
	if err != nil {
		return err
	}
	b.mu.Lock()
	b.consume = consume
	b.mu.Unlock()
	return b.send(msg)
}
