package wizard

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.sr.ht/~tymek/przypominajka/v2/models"
	"git.sr.ht/~tymek/przypominajka/v2/storage"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	addCallbackStepType    = "type"
	addCallbackStepSurname = "surname"
)

const (
	addStepStart int = iota
	addStepMonth
	addStepDay
	addStepType
	addStepName
	addStepSurname
	addStepDone
)

type Add struct {
	id   string
	step int
	done context.CancelFunc
	e    models.Event
}

var _ Interface = (*Add)(nil)

func (a *Add) ID() string {
	return a.id
}

func (a *Add) Active() bool {
	return a.step != addStepStart
}

func (a *Add) Name() string {
	return "add"
}

func (a *Add) start(id string, done context.CancelFunc, update tg.Update) tg.Chattable {
	a.id = id
	a.done = done
	msg, _, _ := a.Next(nil, update)
	return msg
}

var _ Consume = (*Add)(nil).Next

// NOTE: this works on a happy path. If a parsing error occurs, then that's
// likely due to a malicious client using malformed callback. For steps that
// consume update.Message.Text we just retry the same step.
func (a *Add) Next(s storage.Interface, update tg.Update) (tg.Chattable, Consume, error) {
	switch a.step {
	case addStepStart:
		msg := tg.NewMessage(update.FromChat().ID, "Wybierz miesiąc:")
		msg.ReplyMarkup = keyboardMonths(a)
		a.step += 1
		return msg, nil, nil
	case addStepMonth:
		month, err := parseCallbackData(update.CallbackData(), a, callbackPartMonth)
		if err != nil {
			return nil, nil, err
		}
		m, err := strconv.Atoi(month)
		if err != nil {
			return nil, nil, err
		}
		a.e.Month = time.Month(m)
		msg := tg.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, "Wybierz dzień:")
		msg.ReplyMarkup = keyboardDays(a, a.e.Month)
		a.step += 1
		return msg, nil, nil
	case addStepDay:
		day, err := parseCallbackData(update.CallbackData(), a, callbackPartDay)
		if err != nil {
			return nil, nil, err
		}
		d, err := strconv.Atoi(day)
		if err != nil {
			return nil, nil, err
		}
		a.e.Day = d
		msg := tg.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, "Wybierz rodzaj wydarzenia:")
		msg.ReplyMarkup = a.keyboardTypes()
		a.step += 1
		return msg, nil, nil
	case addStepType:
		eventType, err := parseCallbackData(update.CallbackData(), a, addCallbackStepType)
		if err != nil {
			return nil, nil, err
		}
		et := models.EventType(eventType)
		if err := et.Validate(); err != nil {
			return nil, nil, err
		}
		a.e.Type = et
		msg := tg.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, "Wyślij jedno imię lub dwa imiona (każde w osobnej linijce)")
		a.step += 1
		return msg, a.Next, nil
	// FIXME: this step does not run for group chats
	case addStepName:
		lines := strings.Split(strings.TrimSpace(update.Message.Text), "\n")
		if len(lines) != 1 && len(lines) != 2 {
			a.step -= 1
			return tg.NewMessage(update.FromChat().ID, "Wyślij jedno imię lub dwa imiona (każde w osobnej linijce)"), a.Next, nil
		}
		if len(lines) == 1 {
			a.e.Name = lines[0]
		} else {
			a.e.Names = (*[2]string)(lines)
		}
		msg := tg.NewMessage(update.FromChat().ID, "Wyślij nazwisko")
		msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Pomiń", newCallbackData(a, addCallbackStepSurname, "skip")),
			),
		)
		a.step += 1
		return msg, a.Next, nil
	case addStepSurname:
		if update.Message != nil {
			a.e.Surname = update.Message.Text
		}
		if err := a.e.Validate(); err != nil {
			return nil, nil, err
		}
		if err := s.Add(update.FromChat().ID, a.e); err != nil {
			return nil, nil, err
		}
		msg := tg.NewMessage(update.FromChat().ID, fmt.Sprintf("Gotowe! Dodałem:\n%s", a.e.Format(true)))
		a.step += 1
		a.done()
		return msg, nil, nil
	case addStepDone:
		return nil, nil, ErrDone
	}
	return nil, nil, ErrUnknownWizardStep
}

func (a *Add) Reset() {
	a.id = ""
	a.step = addStepStart
	a.e = models.Event{}
	a.done = nil
}

func (a *Add) keyboardTypes() *tg.InlineKeyboardMarkup {
	rows := make([][]tg.InlineKeyboardButton, 1)
	for _, et := range models.EventTypes {
		rows[0] = append(rows[0], tg.NewInlineKeyboardButtonData(et.Format(false), newCallbackData(a, addCallbackStepType, string(et))))
	}
	return &tg.InlineKeyboardMarkup{InlineKeyboard: rows}
}
