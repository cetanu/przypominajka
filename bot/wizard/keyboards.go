package wizard

import (
	"strconv"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	callbackPartMonth = "month"
	callbackPartDay   = "day"
)

func keyboardMonths(w Interface) tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Styczeń", newCallbackData(w, callbackPartMonth, "1")),
			tg.NewInlineKeyboardButtonData("Luty", newCallbackData(w, callbackPartMonth, "2")),
			tg.NewInlineKeyboardButtonData("Marzec", newCallbackData(w, callbackPartMonth, "3")),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Kwiecień", newCallbackData(w, callbackPartMonth, "4")),
			tg.NewInlineKeyboardButtonData("Maj", newCallbackData(w, callbackPartMonth, "5")),
			tg.NewInlineKeyboardButtonData("Czerwiec", newCallbackData(w, callbackPartMonth, "6")),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Lipiec", newCallbackData(w, callbackPartMonth, "7")),
			tg.NewInlineKeyboardButtonData("Sierpień", newCallbackData(w, callbackPartMonth, "8")),
			tg.NewInlineKeyboardButtonData("Wrzesień", newCallbackData(w, callbackPartMonth, "9")),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Październik", newCallbackData(w, callbackPartMonth, "10")),
			tg.NewInlineKeyboardButtonData("Listopad", newCallbackData(w, callbackPartMonth, "11")),
			tg.NewInlineKeyboardButtonData("Grudzień", newCallbackData(w, callbackPartMonth, "12")),
		),
	)
}

func keyboardDays(w Interface, m time.Month) *tg.InlineKeyboardMarkup {
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
		rows[i/nCols] = append(rows[i/nCols], tg.NewInlineKeyboardButtonData(d, newCallbackData(w, callbackPartDay, d)))
	}
	return &tg.InlineKeyboardMarkup{InlineKeyboard: rows}
}
