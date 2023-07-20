package format

const (
	// MessageDone is populated with: caller's username and edited message text.
	MessageDone = "_✅ %s złożył(a) życzenia_\n\n%s"
	MessageOops = "Coś poszło nie tak"
)

// Messages sent out after completing a particular step.
const (
	MessageChooseMonth = "Wybierz miesiąc:"
	MessageChooseDay   = "Wybierz dzień:"
	MessageAddStepDay  = "Wybierz rodzaj wydarzenia:"
	MessageAddStepType = "Wyślij jedno imię lub dwa imiona (każde w osobnej linijce)"
	MessageAddStepName = "Wyślij nazwisko"
)

const (
	MarkupButtonDone = "Gotowe"
	MarkupButtonSkip = "Pomiń"
)
