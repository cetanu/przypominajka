package i18n

import (
	"fmt"
	"strings"
)

var Languages = map[string]map[string]string{
	"pl": {
		"choose_day":          "Wybierz dzień:",
		"choose_event_type":   "Wybierz rodzaj wydarzenia:",
		"choose_month":        "Wybierz miesiąc:",
		"confirm_delete":      "Czy na pewno chcesz usunąć:\n%s?",
		"delete_cancelled":    "Przerwano usuwanie. Wpisz /delete, aby rozpocząć ponownie",
		"deleted":             "Usunięto:\n%s",
		"enter_name":          "Wyślij jedno imię lub dwa imiona (każde w osobnej linijce)",
		"event_added":         "Gotowe! Dodałem:\n%s",
		"language_set":        "Język został ustawiony na %s.",
		"no_events":           "Nie ma żadnych wydarzeń",
		"select_event":        "Wybierz wydarzenie",
		"send_surname":        "Wyślij nazwisko",
		"wishes_sent":         "_✅ %s złożył(a) życzenia_\n\n%s",
		"done":                "Gotowe",
		"no_events_on_day":    "W wybranym dniu nie ma żadnych wydarzeń. Wpisz /delete, aby rozpocząć ponownie",
		"operation_cancelled": "Przerwano!",
		"skip":                "Pomiń",
		"something_wrong":     "Coś poszło nie tak",
	},
	"en": {
		"choose_day":          "Choose a day:",
		"choose_event_type":   "Select the type of event:",
		"choose_month":        "Choose a month:",
		"confirm_delete":      "Are you sure you want to delete:\n%s?",
		"delete_cancelled":    "Deletion canceled. Type /delete to restart.",
		"deleted":             "Deleted:\n%s",
		"enter_name":          "Send one name or two names (each on a separate line)",
		"event_added":         "Done! Added:\n%s",
		"language_set":        "Language set to %s.",
		"no_events":           "There are no events",
		"no_events_on_day":    "There are no events on the selected day. Type /delete to restart.",
		"select_event":        "Select an event",
		"send_surname":        "Send surname",
		"wishes_sent":         "_✅ %s sent wishes_\n\n%s",
		"done":                "Done",
		"operation_cancelled": "Cancelled!",
		"skip":                "Skip",
		"something_wrong":     "Something went wrong",
	},
}

// Usage:
// lang := GetUserLanguage(update.FromChat().ID) // Fetch user language
// msg := i18n.T(lang, "choose_month") // Fetch translated message
func T(lang, key string, args ...interface{}) string {
	if translation, ok := Languages[lang][key]; ok {
		return fmt.Sprintf(translation, args...)
	}
	return fmt.Sprintf("Missing translation: %s", key)
}

func GetSupportedLanguages(delimiter string) string {
	keys := make([]string, 0, len(Languages))
	for lang := range Languages {
		keys = append(keys, lang)
	}
	return strings.Join(keys, delimiter)
}
