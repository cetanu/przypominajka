package format

// Event type formats
const (
	Birthday           = "urodziny"
	Nameday            = "imieniny"
	WeddingAnniversary = "rocznicę ślubu"
)

const (
	// Singular and ListSingular are populated with: name and event type format.
	Singular     = "%s obchodzi dziś %s!"
	ListSingular = "%s obchodzi %s"
	// Singular and ListSingular are populated with: name, surname, and event type format.
	SingularSurname     = "%s %s obchodzi dziś %s!"
	ListSingularSurname = "%s %s obchodzi %s"
	// MessagePlural and ListMessagePlural are populated with: names[0], names[1], and event type format.
	MessagePlural     = "%s i %s obchodzą dziś %s!"
	ListMessagePlural = "%s i %s obchodzą %s"
	// MessagePluralSurname and ListMessagePluralSurname are populated with: names[0], names[1], surname, and event type format.
	MessagePluralSurname     = "%s i %s %s obchodzą dziś %s!"
	ListMessagePluralSurname = "%s i %s %s obchodzą %s"
	// ListLine is populated with: day, month, and event.
	ListLine = "%02d.%02d - %s"
)

const (
	// MsgDone is populated with: caller's username and edited message text.
	MsgDone     = "_✅ %s złożył(a) życzenia_\n\n%s"
	MsgNoEvents = "Nie znalazłem żadnych wydarzeń"
)
