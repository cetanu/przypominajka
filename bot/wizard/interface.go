package wizard

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"git.sr.ht/~tymek/przypominajka/storage"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const CallbackSep = ":"

var (
	ErrDone                = errors.New("wizard already done")
	ErrInvalidCallbackData = errors.New("invalid callback data")
	ErrUnknownWizardStep   = errors.New("unknown wizard step")
)

type Interface interface {
	ID() string
	Active() bool
	Name() string
	Next(s storage.Interface, update tg.Update) (tg.Chattable, Consume, error)
	Reset()
	start(id string, done context.CancelFunc, update tg.Update) tg.Chattable
}

type Consume func(s storage.Interface, update tg.Update) (tg.Chattable, Consume, error)

func Start(w Interface, update tg.Update) tg.Chattable {
	w.Reset()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-ctx.Done():
		case <-time.After(5 * time.Minute):
			w.Reset()
		}
	}()
	return w.start(newID(), cancel, update)
}

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func newID() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func newCallbackData(w Interface, parts ...string) string {
	return fmt.Sprint(w.Name(), CallbackSep, w.ID(), CallbackSep, strings.Join(parts, CallbackSep))
}

func parseCallbackData(s string, w Interface, static ...string) (string, error) {
	parts := strings.Split(s, CallbackSep)
	if len(parts) != len(static)+3 || parts[0] != w.Name() || parts[1] != w.ID() {
		return "", ErrInvalidCallbackData
	}
	for i, p := range static {
		if p != parts[i+2] {
			return "", ErrInvalidCallbackData
		}
	}
	return parts[len(parts)-1], nil
}
