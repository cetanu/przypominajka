package main

import "fmt"

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
		return formatBirthday
	case nameday:
		return formatNameday
	case wedding:
		return formatWeddingAnniversary
	}
	return string(et)
}
