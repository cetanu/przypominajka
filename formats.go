package main

// Event type formats
const (
	formatBirthday           = "urodziny"
	formatNameday            = "imieniny"
	formatWeddingAnniversary = "rocznicę ślubu"
)

const (
	// formatSingular and formatListSingular are populated with: name and event type format.
	formatSingular     = "%s obchodzi dziś %s!"
	formatListSingular = "%s obchodzi %s"
	// formatSingular and formatListSingular are populated with: name, surname, and event type format.
	formatSingularSurname     = "%s %s obchodzi dziś %s!"
	formatListSingularSurname = "%s %s obchodzi %s"
	// formatMessagePlural and formatListMessagePlural are populated with: names[0], names[1], and event type format.
	formatMessagePlural     = "%s i %s obchodzą dziś %s!"
	formatListMessagePlural = "%s i %s obchodzą %s"
	// formatMessagePluralSurname and formatListMessagePluralSurname are populated with: names[0], names[1], surname, and event type format.
	formatMessagePluralSurname     = "%s i %s %s obchodzą dziś %s!"
	formatListMessagePluralSurname = "%s i %s %s obchodzą %s"
	// formatDone is populated with: caller's username and edited message text.
	formatDone = "_✅ %s złożył(a) życzenia_\n\n%s"
	// formatListLine is populated with: day, month, and event.
	formatListLine = "%02d.%02d - %s\n"
)
