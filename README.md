# przypominajka

przypominajka runs a Telegram bot to send notifications about birthdays,
namedays, and anniversaries.

Note: all user-facing messages are in Polish. If you would like to translate
them, or introduce internationalization, then feel free to reach out to me, I
will gladly help.

## Installation
Run `make` to compile przypominajka.
Run `make install` to install przypominajka and completions to `/usr/local/`.
Clean up with `make clean` and `make uninstall`, respectively.

To override `/usr/local/` PREFIX variable use `make -e PREFIX=/foo/bar/baz/`.

## Usage
```
przypominajka - a Telegram bot for sending event reminders

Description:
  przypominajka reads a YAML file with events and sends reminders about them.
  The reminders are sent out on the day of the event between 08:30 and 09:29
  system time (exact time depends on serve command startup time).

  All user-facing messages are written in Polish.

Example events.yaml:
  chat_ids:
    - 1234
    - 5678 # has no initial data

  data:
    1234: # Chat ID 1234
      1: # January
        5:
          - name: "John"
            type: "birthday"
          - name: "Jane"
            surname: "Doe"
            type: "nameday"
      4: # April
        17:
          - names: ["John", "Jane"]
            surname: "Doe"
            type: "wedding anniversary"

Notes:
  - Name and names are mutually exclusive.
  - Names, if present, must have two elements.
  - Surname is optional.
  - Type has to be one of: birthday, nameday, wedding anniversary.

Usage:
  przypominajka [command]

Available Commands:
  bot         Start Telegram bot to serve events and listen for updates
  help        Help about any command
  show        Show events

Flags:
      --events string   YAML file defining events (default "events.yaml")
  -h, --help            help for przypominajka
  -v, --version         version for przypominajka

Use "przypominajka [command] --help" for more information about a command.
```

### Bot
```
Start Telegram bot to serve events and listen for updates

Usage:
  przypominajka bot [flags]

Flags:
  -h, --help           help for bot
      --token string   Telegram bot token

Global Flags:
      --events string   YAML file defining events (default "events.yaml")
```

### Show
```
Show events

Usage:
  przypominajka show [command]

Available Commands:
  all         Show all events
  next        Find the next day with events and list them

Flags:
      --chat-id int   Telegram chat ID
  -h, --help          help for show

Global Flags:
      --events string   YAML file defining events (default "events.yaml")

Use "przypominajka show [command] --help" for more information about a command.
```

####  All
```
Show all events

Usage:
  przypominajka show all [flags]

Flags:
  -h, --help   help for all

Global Flags:
      --chat-id int     Telegram chat ID
      --events string   YAML file defining events (default "events.yaml")
```

####  Next
```
Find the next day with events and list them

Usage:
  przypominajka show next [flags]

Flags:
  -h, --help   help for next

Global Flags:
      --chat-id int     Telegram chat ID
      --events string   YAML file defining events (default "events.yaml")
```

## Bot Setup
Define the following commands for the bot:
```
abort - Przerwij dodawanie lub usuwanie
add - Dodaj nowe wydarzenie
delete - Usuń wydarzenie
list - Wypisz wszystkie wydarzenia
next - Pokaż następne wydarzenia
```
