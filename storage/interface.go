package storage

import (
	"errors"
	"time"

	"git.sr.ht/~tymek/przypominajka/models"
)

type Interface interface {
	Format(chatID int64) string
	ChatIDs() []int64
	IsEnabled(chatID int64) bool

	At(chatID int64, t time.Time) (models.Events, error)
	Add(chatID int64, e models.Event) error
	Remove(chatID int64, e models.Event) error
}

func At(s Interface, chatID int64, m time.Month, d int) (models.Events, error) {
	return s.At(chatID, time.Date(1970, m, d, 0, 0, 0, 0, time.UTC))
}

func Next(s Interface, chatID int64) (models.Events, error) {
	now := time.Now()
	nextDay := now.AddDate(0, 0, 1)
	nextYear := nextDay.AddDate(1, 0, 0)
	for t := nextDay; t.Before(nextYear); t = t.AddDate(0, 0, 1) {
		events, err := s.At(chatID, t)
		if err != nil {
			return nil, err
		}
		if len(events) > 0 {
			return events, nil
		}
	}
	return nil, errors.New("Nie ma żadnych wydarzeń")
}

func Today(s Interface, chatID int64) (models.Events, error) {
	return s.At(chatID, time.Now())
}
