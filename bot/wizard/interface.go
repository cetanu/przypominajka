package wizard

import (
	"errors"
	"time"

	"git.sr.ht/~tymek/przypominajka/storage"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var ErrDone = errors.New("wizard: done")

type Interface interface {
	Name() string
	Start() tg.Chattable
	Next(s storage.Interface, update tg.Update) (tg.Chattable, error)
	Reset()
}

func ResetAfter(w Interface, d time.Duration) {
	time.AfterFunc(d, w.Reset)
}
