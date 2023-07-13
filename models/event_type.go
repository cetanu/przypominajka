package models

import (
	"fmt"

	"git.sr.ht/~tymek/przypominajka/format"
)

const (
	Birthday EventType = "birthday"
	Nameday  EventType = "nameday"
	Wedding  EventType = "wedding anniversary"
)

type EventType string

var _ fmt.Stringer = EventType("")

func (et EventType) String() string {
	switch et {
	case Birthday:
		return format.Birthday
	case Nameday:
		return format.Nameday
	case Wedding:
		return format.WeddingAnniversary
	}
	return string(et)
}
