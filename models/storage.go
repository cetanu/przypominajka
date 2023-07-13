package models

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type Storage interface {
	At(t time.Time) (Events, error)
}

func Next(s Storage) (Events, error) {
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
	return nil, ErrNotFound
}
