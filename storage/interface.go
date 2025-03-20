package storage

import (
	"errors"
	"time"

	"github.com/TymekDev/przypominajka/v2/i18n"
	"github.com/TymekDev/przypominajka/v2/models"
)

type Interface interface {
	Format(chatID int64) string
	ChatIDs() []int64
	IsEnabled(chatID int64) bool

	At(chatID int64, t time.Time) (models.Events, error)
	Add(chatID int64, e models.Event) error
	Remove(chatID int64, e models.Event) error
	GetUserLanguage(chatID int64) string
	SetUserLanguage(chatID int64, lang string) error
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
	lang := s.GetUserLanguage(chatID)
	return nil, errors.New(i18n.T(lang, "no_events"))
}

func Today(s Interface, chatID int64) (models.Events, error) {
	return s.At(chatID, time.Now())
}
