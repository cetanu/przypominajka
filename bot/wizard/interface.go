package wizard

import (
	"errors"
	"fmt"
	"strings"

	"git.sr.ht/~tymek/przypominajka/storage"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const CallbackSep = ":"

var (
	ErrDone                = errors.New("wizard already done")
	ErrInvalidCallbackData = errors.New("invalid callback data")
	ErrUserError           = errors.New("user error")
)

type Interface interface {
	Active() bool
	Name() string
	Start(update tg.Update) tg.Chattable
	Next(s storage.Interface, update tg.Update) (tg.Chattable, Consume, error)
	Reset()
}

type Consume func(s storage.Interface, update tg.Update) (tg.Chattable, Consume, error)

func newCallbackData(w Interface, parts ...string) string {
	return fmt.Sprint(w.Name(), CallbackSep, strings.Join(parts, CallbackSep))
}

func parseCallbackData(s string, w Interface, static ...string) (string, error) {
	parts := strings.Split(s, CallbackSep)
	if len(parts) != len(static)+2 || parts[0] != w.Name() {
		return "", ErrInvalidCallbackData
	}
	for i, p := range static {
		if p != parts[i+1] {
			return "", ErrInvalidCallbackData
		}
	}
	return parts[len(parts)-1], nil
}
