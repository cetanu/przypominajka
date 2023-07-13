package models

import (
	"time"
)

type Storage interface {
	At(t time.Time) (Events, error)
}
