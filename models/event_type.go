package models

import (
	"fmt"

	"git.sr.ht/~tymek/przypominajka/format"
)

const (
	birthday eventType = "birthday"
	nameday  eventType = "nameday"
	wedding  eventType = "wedding anniversary"
)

type eventType string

var _ fmt.Stringer = eventType("")

func (et eventType) String() string {
	switch et {
	case birthday:
		return format.Birthday
	case nameday:
		return format.Nameday
	case wedding:
		return format.WeddingAnniversary
	}
	return string(et)
}
