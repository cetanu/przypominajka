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

func At(s Interface, m time.Month, d int) (models.Events, error) {
	return s.At(time.Date(1970, m, d, 0, 0, 0, 0, time.UTC))
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
