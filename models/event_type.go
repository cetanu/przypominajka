package models

import (
	"fmt"
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
		return formatBirthday
	case Nameday:
		return formatNameday
	case Wedding:
		return formatWeddingAnniversary
	}
	return string(et)
}
