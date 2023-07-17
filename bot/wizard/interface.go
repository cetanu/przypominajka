package wizard

import (
	"errors"
	"fmt"
	"strings"

	"git.sr.ht/~tymek/przypominajka/storage"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const CallbackSep = ":"

var ErrDone = errors.New("wizard done")

type Interface interface {
	Name() string
	Start(update tg.Update) tg.Chattable
	Next(s storage.Interface, update tg.Update) (tg.Chattable, error)
	Reset()
}

func newCallbackData(w Interface, parts ...string) string {
	return fmt.Sprint(w.Name(), CallbackSep, strings.Join(parts, CallbackSep))
}
