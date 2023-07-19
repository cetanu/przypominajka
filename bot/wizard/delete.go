package wizard

import (
	"context"

	"git.sr.ht/~tymek/przypominajka/storage"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Delete struct {
	id   string
	done context.CancelFunc
}

var _ Interface = (*Delete)(nil)

func (d *Delete) ID() string {
	return d.id
}

func (d *Delete) Name() string {
	return "delete"
}

func (d *Delete) Active() bool {
	return false
}

func (d *Delete) start(id string, done context.CancelFunc, update tg.Update) tg.Chattable {
	d.id = id
	d.done = done
	return nil
}

func (d *Delete) Next(s storage.Interface, update tg.Update) (tg.Chattable, Consume, error) {
	return nil, nil, ErrUnknownWizardStep
}

func (d *Delete) Reset() {
	d.id = ""
	d.done = nil
}
