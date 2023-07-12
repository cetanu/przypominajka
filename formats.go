package main

// Event type formats
const (
	formatBirthday           = "urodziny"
	formatNameday            = "imieniny"
	formatWeddingAnniversary = "rocznicę ślubu"
)

const (
	// formatSingular is populated with: name and event type format.
	formatSingular = "%s ma dziś %s!"
	// formatSingular is populated with: name, surname, and event type format.
	formatSingularSurname = "%s %s ma dziś %s!"
	// formatMessagePlural is populated with: names[0], names[1], and event type format.
	formatMessagePlural = "%s i %s mają dziś %s!"
	// formatMessagePluralSurname is populated with: names[0], names[1], surname, and event type format.
	formatMessagePluralSurname = "%s i %s %s mają dziś %s!"
	// formatDone is populated with: caller's username and edited message text.
	formatDone = "_✅ %s złożył(a) życzenia_\n\n%s"
	// formatListLine is populated with: day, month, and event.
	formatListLine = "%02d.%02d - %s\n"
)
