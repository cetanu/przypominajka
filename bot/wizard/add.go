package wizard

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.sr.ht/~tymek/przypominajka/format"
	"git.sr.ht/~tymek/przypominajka/models"
	"git.sr.ht/~tymek/przypominajka/storage"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	addCallbackStepMonth   = "month"
	addCallbackStepDay     = "day"
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
	e    models.Event
	done context.CancelFunc
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
		msg := tg.NewMessage(update.FromChat().ID, format.MessageAddStepStart)
		msg.ReplyMarkup = a.keyboardMonths()
		a.step += 1
		return msg, nil, nil
	case addStepMonth:
		month, err := parseCallbackData(update.CallbackData(), a, addCallbackStepMonth)
		if err != nil {
			return nil, nil, err
		}
		m, err := strconv.Atoi(month)
		if err != nil {
			return nil, nil, err
		}
		a.e.Month = time.Month(m)
		msg := tg.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, format.MessageAddStepMonth)
		msg.ReplyMarkup = a.keyboardDays(a.e.Month)
		a.step += 1
		return msg, nil, nil
	case addStepDay:
		day, err := parseCallbackData(update.CallbackData(), a, addCallbackStepDay)
		if err != nil {
			return nil, nil, err
		}
		d, err := strconv.Atoi(day)
		if err != nil {
			return nil, nil, err
		}
		a.e.Day = d
		msg := tg.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, format.MessageAddStepDay)
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
		msg := tg.NewEditMessageText(update.FromChat().ID, update.CallbackQuery.Message.MessageID, format.MessageAddStepType)
		a.step += 1
		return msg, a.Next, nil
	case addStepName:
		lines := strings.Split(strings.TrimSpace(update.Message.Text), "\n")
		if len(lines) != 1 && len(lines) != 2 {
			a.step -= 1
			return tg.NewMessage(update.FromChat().ID, format.MessageAddStepType), a.Next, nil
		}
		if len(lines) == 1 {
			a.e.Name = lines[0]
		} else {
			a.e.Names = (*[2]string)(lines)
		}
		msg := tg.NewMessage(update.FromChat().ID, format.MessageAddStepName)
		msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData(format.MarkupButtonSkip, newCallbackData(a, addCallbackStepSurname, "skip")),
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
		if err := s.Add(a.e); err != nil {
			return nil, nil, err
		}
		msg := tg.NewMessage(update.FromChat().ID, fmt.Sprintf("Gotowe! Dodałem:\n%s", a.e.Format(true)))
		a.step += 1
		a.done()
		return msg, nil, nil
	case addStepDone:
		return nil, nil, ErrDone
	}
	return nil, nil, errors.New("unknown wizard step")
}

func (a *Add) Reset() {
	a.id = ""
	a.step = addStepStart
	a.e = models.Event{}
	a.done = nil
}

func (a *Add) keyboardMonths() tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Styczeń", a.callbackMonth(1)),
			tg.NewInlineKeyboardButtonData("Luty", a.callbackMonth(2)),
			tg.NewInlineKeyboardButtonData("Marzec", a.callbackMonth(3)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Kwiecień", a.callbackMonth(4)),
			tg.NewInlineKeyboardButtonData("Maj", a.callbackMonth(5)),
			tg.NewInlineKeyboardButtonData("Czerwiec", a.callbackMonth(6)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Lipiec", a.callbackMonth(7)),
			tg.NewInlineKeyboardButtonData("Sierpień", a.callbackMonth(8)),
			tg.NewInlineKeyboardButtonData("Wrzesień", a.callbackMonth(9)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Październik", a.callbackMonth(10)),
			tg.NewInlineKeyboardButtonData("Listopad", a.callbackMonth(11)),
			tg.NewInlineKeyboardButtonData("Grudzień", a.callbackMonth(12)),
		),
	)
}

func (a *Add) callbackMonth(m int) string {
	return newCallbackData(a, addCallbackStepMonth, strconv.Itoa(m))
}

func (a *Add) keyboardDays(m time.Month) *tg.InlineKeyboardMarkup {
	const nCols = 8 // that's the max Telegram allows for an inline keyboard
	n := 31
	switch m {
	case time.February:
		n = 29
	case time.April, time.June, time.September, time.November:
		n = 30
	}
	rows := make([][]tg.InlineKeyboardButton, 4)
	for i := 0; i < n; i++ {
		d := strconv.Itoa(i + 1)
		rows[i/nCols] = append(rows[i/nCols], tg.NewInlineKeyboardButtonData(d, newCallbackData(a, addCallbackStepDay, d)))
	}
	return &tg.InlineKeyboardMarkup{InlineKeyboard: rows}
}

func (a *Add) keyboardTypes() *tg.InlineKeyboardMarkup {
	rows := make([][]tg.InlineKeyboardButton, 1)
	for _, et := range models.EventTypes {
		rows[0] = append(rows[0], tg.NewInlineKeyboardButtonData(et.Format(false), newCallbackData(a, addCallbackStepType, string(et))))
	}
	return &tg.InlineKeyboardMarkup{InlineKeyboard: rows}
}
