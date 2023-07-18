package format

const (
	// MessageDone is populated with: caller's username and edited message text.
	MessageDone = "_✅ %s złożył(a) życzenia_\n\n%s"
)

// Messages sent out after completing a particular step.
// Therefore MessageAddStepStart mentions selecting months that are the next step.
const (
	MessageAddStepStart = "Wybierz miesiąc:"
	MessageAddStepMonth = "Wybierz dzień:"
	MessageAddStepDay   = "Wybierz rodzaj wydarzenia:"
	MessageAddStepType  = "Wyślij jedno imię lub dwa imiona (każde w osobnej linijce)"
	MessageAddStepName  = "Wyślij nazwisko"
)

const (
	MarkupButtonDone = "Gotowe"
	MarkupButtonSkip = "Pomiń"
)
