package models

import (
	"fmt"

	"git.sr.ht/~tymek/przypominajka/format"
)

var EventTypes = []EventType{Birthday, Nameday, Wedding}

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

func (et EventType) Validate() error {
	for _, t := range EventTypes {
		if et == t {
			return nil
		}
	}
	return fmt.Errorf("invalid event type: %s", et)
}
