package wizard

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"git.sr.ht/~tymek/przypominajka/v2/models"
	"git.sr.ht/~tymek/przypominajka/v2/storage"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	deleteCallbackStepEvent = "event"
)

const (
	deleteStepStart int = iota
	deleteStepMonth
	deleteStepDay
	deleteStepConfirm
	deleteStepRemove
	deleteStepDone
)

type Delete struct {
	id    string
	step  int
	done  context.CancelFunc
	month time.Month
	day   int
	e     models.Event
}

var _ Interface = (*Delete)(nil)

func (d *Delete) ID() string {
	return d.id
}

func (d *Delete) Name() string {
	return "delete"
}

func (d *Delete) Active() bool {
	return d.step != deleteStepStart
}

func (d *Delete) start(id string, done context.CancelFunc, update tg.Update) tg.Chattable {
	d.id = id
	d.done = done
	msg, _, _ := d.Next(nil, update)
	return msg
}

func (d *Delete) Next(s storage.Interface, update tg.Update) (tg.Chattable, Consume, error) {
	switch d.step {
	case deleteStepStart:
		msg := tg.NewMessage(update.FromChat().ID, "Wybierz miesiąc:")
		msg.ReplyMarkup = keyboardMonths(d)
		d.step += 1
		return msg, nil, nil
	case deleteStepMonth:
		month, err := parseCallbackData(update.CallbackData(), d, callbackPartMonth)
		if err != nil {
			return nil, nil, err
		}
		m, err := strconv.Atoi(month)
		if err != nil {
			return nil, nil, err
		}
		d.month = time.Month(m)
		msg := tg.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, "Wybierz dzień:")
		msg.ReplyMarkup = keyboardDays(d, d.month)
		d.step += 1
		return msg, nil, nil
	case deleteStepDay:
		day, err := parseCallbackData(update.CallbackData(), d, callbackPartDay)
		if err != nil {
			return nil, nil, err
		}
		if d.day, err = strconv.Atoi(day); err != nil {
			return nil, nil, err
		}
		events, err := storage.At(s, update.FromChat().ID, d.month, d.day)
		if err != nil {
			return nil, nil, err
		}
		if len(events) == 0 {
			msg := tg.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, "W wybranym dniu nie ma żadnych wydarzeń. Wpisz /delete, aby rozpocząć ponownie")
			return msg, nil, nil
		}
		msg := tg.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, "Wybierz wydarzenie")
		msg.ReplyMarkup = d.keyboardEvents(events)
		d.step += 1
		return msg, nil, nil
	case deleteStepConfirm:
		index, err := parseCallbackData(update.CallbackData(), d, deleteCallbackStepEvent)
		if err != nil {
			return nil, nil, err
		}
		i, err := strconv.Atoi(index)
		if err != nil {
			return nil, nil, err
		}
		events, err := storage.At(s, update.FromChat().ID, d.month, d.day)
		if err != nil {
			return nil, nil, err
		}
		if i >= len(events) {
			return nil, nil, errors.New("something changed and things broke")
		}
		d.e = events[i]
		msg := tg.NewEditMessageTextAndMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, fmt.Sprintf("Czy na pewno chcesz usunąć:\n%s?", d.e.Format(true)),
			tg.NewInlineKeyboardMarkup(
				tg.NewInlineKeyboardRow(
					tg.NewInlineKeyboardButtonData("Tak", newCallbackData(d, "confirm", "yes")),
					tg.NewInlineKeyboardButtonData("Nie", newCallbackData(d, "confirm", "no")),
				),
			))
		d.step += 1
		return msg, nil, nil
	case deleteStepRemove:
		value, err := parseCallbackData(update.CallbackData(), d, "confirm")
		if err != nil {
			return nil, nil, err
		}
		text := "Przerwano usuwanie. Wpisz /delete, aby rozpocząć ponownie"
		if value == "yes" {
			if err := s.Remove(update.FromChat().ID, d.e); err != nil {
				return nil, nil, err
			}
			text = fmt.Sprintf("Usunięto:\n%s", d.e.Format(true))
		}
		d.step += 1
		d.done()
		msg := tg.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, text)
		return msg, nil, nil
	case deleteStepDone:
		return nil, nil, ErrDone
	}
	return nil, nil, ErrUnknownWizardStep
}

func (d *Delete) Reset() {
	d.id = ""
	d.step = addStepStart
	d.done = nil
	d.month = 0
	d.day = 0
	d.e = models.Event{}
}

func (d *Delete) keyboardEvents(events models.Events) *tg.InlineKeyboardMarkup {
	rows := make([][]tg.InlineKeyboardButton, len(events))
	for i, e := range events {
		rows[i] = []tg.InlineKeyboardButton{tg.NewInlineKeyboardButtonData(e.Format(true), newCallbackData(d, deleteCallbackStepEvent, strconv.Itoa(i)))}
	}
	return &tg.InlineKeyboardMarkup{InlineKeyboard: rows}
}
