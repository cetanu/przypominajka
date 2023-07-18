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
	return et.Format(true)
}

func (et EventType) Format(accusative bool) string {
	switch {
	case et == Birthday && accusative:
		return format.BirthdayAccusative
	case et == Birthday:
		return format.BirthdayNominative
	case et == Nameday && accusative:
		return format.NamedayAccusative
	case et == Nameday:
		return format.NamedayNominative
	case et == Wedding && accusative:
		return format.WeddingAnniversaryAccusative
	case et == Wedding:
		return format.WeddingAnniversaryNominative
	}
	return string(et)
}

func (et EventType) Validate() error {
	for _, t := range EventTypes {
		if et == t {
			return nil
		}
	}
	return ErrInvalidEventType
}
