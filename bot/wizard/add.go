package wizard

import (
	"context"
	"errors"
	"time"

	"git.sr.ht/~tymek/przypominajka/models"
	"git.sr.ht/~tymek/przypominajka/storage"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	step        int
	e           models.Event
	cancelReset context.CancelFunc
}

var _ Interface = (*Add)(nil)

func (a *Add) Name() string {
	return "add"
}

func (a *Add) Start() tg.Chattable {
	ctx, cancel := context.WithCancel(context.Background())
	a.cancelReset = cancel
	go func() {
		select {
		case <-ctx.Done():
		case <-time.After(30 * time.Second):
			a.Reset()
		}
	}()
	return nil
}

func (a *Add) Next(s storage.Interface, update tg.Update) (tg.Chattable, error) {
	switch a.step {
	case addStepMonth:
	case addStepDay:
	case addStepType:
	case addStepName:
	case addStepSurname:
	case addStepDone:
		return nil, ErrDone
	}
	a.step += 1
	return nil, errors.New("unknown wizard step")
}

func (a *Add) Reset() {
	if cr := a.cancelReset; cr != nil {
		cr()
	}
	a.step = addStepStart
	a.e = models.Event{}
	a.cancelReset = nil
}
