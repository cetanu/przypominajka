package storage

import (
	"fmt"
	"time"

	"git.sr.ht/~tymek/przypominajka/models"
)

type Interface interface {
	fmt.Stringer
	At(t time.Time) (models.Events, error)
	Add(e models.Event) error
	Remove(e models.Event) error
}

func Next(s Interface) (models.Events, error) {
	now := time.Now()
	nextDay := now.AddDate(0, 0, 1)
	nextYear := nextDay.AddDate(1, 0, 0)
	for t := nextDay; t.Before(nextYear); t = t.AddDate(0, 0, 1) {
		events, err := s.At(t)
		if err != nil {
			return nil, err
		}
		if len(events) > 0 {
			return events, nil
		}
	}
	return nil, models.ErrNotFound
}

func Today(s Interface) (models.Events, error) {
	return s.At(time.Now())
}
