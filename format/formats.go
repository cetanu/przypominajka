package format

// Events formats
const (
	NoEvents = "Nie ma żadnych wydarzeń"
)

// Event formats
const (
	// Singular and ListSingular are populated with: name and event type.
	Singular     = "%s obchodzi dziś %s!"
	ListSingular = "%s obchodzi %s"
	// SingularSurname and ListSingularSurname are populated with: name, surname, and event type.
	SingularSurname     = "%s %s obchodzi dziś %s!"
	ListSingularSurname = "%s %s obchodzi %s"
	// Plural and ListPlural are populated with: names[0], names[1], and event type.
	Plural     = "%s i %s obchodzą dziś %s!"
	ListPlural = "%s i %s obchodzą %s"
	// PluralSurname and ListPluralSurname are populated with: names[0], names[1], surname, and event type.
	PluralSurname     = "%s i %s %s obchodzą dziś %s!"
	ListPluralSurname = "%s i %s %s obchodzą %s"
	// ListLine is populated with: day, month, and event.
	ListLine = "%02d.%02d - %s"
)

// EventType formats
const (
	BirthdayNominative           = "urodziny"
	BirthdayAccusative           = "urodziny"
	NamedayNominative            = "imieniny"
	NamedayAccusative            = "imieniny"
	WeddingAnniversaryNominative = "rocznica ślubu"
	WeddingAnniversaryAccusative = "rocznicę ślubu"
)
