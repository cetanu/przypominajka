package models

import (
	"fmt"
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
	case et == Birthday:
		return "urodziny"
	case et == Nameday:
		return "imieniny"
	case et == Wedding && accusative:
		return "rocznicę ślubu"
	case et == Wedding:
		return "rocznica ślubu"
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
