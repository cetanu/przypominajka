package bot

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/TymekDev/przypominajka/v2/bot/wizard"
	"github.com/TymekDev/przypominajka/v2/models"
	"github.com/TymekDev/przypominajka/v2/storage"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const dataNotifyDone = "done"

type Bot struct {
	mu      sync.Mutex
	api     *tg.BotAPI
	s       storage.Interface
	wizards map[int64]map[string]wizard.Interface
	consume map[int64]wizard.Consume
}

func New(token string, s storage.Interface) (*Bot, error) {
	api, err := tg.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	log.Println("INFO", "Authorized as", api.Self.UserName)

	m := map[int64]map[string]wizard.Interface{}
	for _, chatID := range s.ChatIDs() {
		m[chatID] = map[string]wizard.Interface{
			(*wizard.Add)(nil).Name():    &wizard.Add{},
			(*wizard.Delete)(nil).Name(): &wizard.Delete{},
		}
	}
	return &Bot{
		api:     api,
		s:       s,
		wizards: m,
		consume: map[int64]wizard.Consume{},
	}, nil
}

func ListenAndServe(token string, s storage.Interface) error {
	b, err := New(token, s)
	if err != nil {
		return err
	}
	go b.Listen()
	b.Serve()
	return nil
}

func (b *Bot) Notify(chatID int64, e models.Event) error {
	msg := tg.NewMessage(chatID, e.Format(false))
	msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Gotowe", dataNotifyDone),
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
				if err := b.send(tg.NewMessage(update.FromChat().ID, "Coś poszło nie tak")); err != nil {
					log.Println("ERROR", "couldn't send internal error message:", err)
				}
			}
		}(update)
	}
}

func (b *Bot) Serve() {
	for t := range time.Tick(time.Hour) {
		if t.UTC().Round(time.Hour).Hour() != 7 { // run once a day between 7:30 and 8:29 UTC
			continue
		}
		for _, chatID := range b.s.ChatIDs() {
			go func(chatID int64) {
				events, err := storage.Today(b.s, chatID)
				if err != nil {
					log.Println("ERROR", err)
					return
				}
				for _, e := range events {
					if err := b.Notify(chatID, e); err != nil {
						log.Println("ERROR", err)
					}
				}
			}(chatID)
		}
	}
}

func (b *Bot) handle(update tg.Update) error {
	if chat := update.FromChat(); chat == nil || !b.s.IsEnabled(chat.ID) {
		return nil
	}
	switch {
	case update.CallbackQuery != nil:
		if _, err := b.api.Request(tg.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)); err != nil {
			return err
		}
		switch data := update.CallbackQuery.Data; data {
		case dataNotifyDone:
			cq := update.CallbackQuery
			const format = "_✅ %s złożył(a) życzenia_\n\n%s"
			edit := tg.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, fmt.Sprintf(format, cq.From.UserName, cq.Message.Text))
			edit.ParseMode = tg.ModeMarkdown
			return b.send(edit)
		default:
			name, _, ok := strings.Cut(update.CallbackData(), wizard.CallbackSep)
			if !ok {
				return nil
			}
			w, ok := b.wizards[update.FromChat().ID][name]
			if !ok {
				return nil
			}
			if w.Active() {
				return b.runConsume(w.Next, update)
			}
			return nil
		}

	case update.Message != nil:
		switch m := update.Message; {
		case m.IsCommand():
			// NOTE: if another bot has /next command, then this will be triggered.
			// To prevent this behavior, we can CommandWithAt() and check whether
			// <command>@<bot_name> matches.
			switch cmd := update.Message.Command(); cmd {
			case "abort":
				b.mu.Lock()
				for _, w := range b.wizards[update.FromChat().ID] {
					w.Reset()
				}
				b.mu.Unlock()
				return b.send(tg.NewMessage(update.FromChat().ID, "Przerwano!"))
			case "list":
				return b.send(tg.NewMessage(update.FromChat().ID, b.s.Format(update.FromChat().ID)))
			case "next":
				events, err := storage.Next(b.s, update.FromChat().ID)
				if err != nil {
					return err
				}
				return b.send(tg.NewMessage(update.FromChat().ID, events.String()))
			default:
				if w, ok := b.wizards[update.FromChat().ID][update.Message.Command()]; ok {
					return b.send(wizard.Start(w, update))
				}
			}

		case m.Text != "":
			return b.runConsume(b.consume[update.FromChat().ID], update)
		}
	}
	return nil
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
	b.mu.Lock()
	msg, consume, err := c(b.s, update)
	b.mu.Unlock()
	if err != nil {
		return err
	}
	b.mu.Lock()
	b.consume[update.FromChat().ID] = consume
	b.mu.Unlock()
	return b.send(msg)
}
