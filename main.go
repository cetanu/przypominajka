package main

import (
	"fmt"
	"log"

	"git.sr.ht/~tymek/przypominajka/bot"
	"git.sr.ht/~tymek/przypominajka/storage"
	"github.com/spf13/cobra"
)

var version string

const description = `przypominajka - a Telegram bot for sending event reminders

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
  - Type has to be one of: birthday, nameday, wedding anniversary.`

func main() {
	log.SetFlags(0)

	var (
		s          storage.Interface
		eventsPath string
		token      string
		chatID     int64
	)

	cmd := &cobra.Command{
		Use:               "przypominajka",
		Short:             description,
		Version:           version,
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Parent() != nil && (cmd.Parent().Use == "completion" || cmd.Parent().Use == "help") {
				return nil
			}
			y, err := storage.NewYAML(eventsPath)
			if err != nil {
				return err
			}
			s = y
			return nil
		},
	}

	cmdShow := &cobra.Command{
		Use:   "show",
		Short: "Show events",
	}

	cmdShowAll := &cobra.Command{
		Use:   "all",
		Short: "Show all events",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(s.Format(chatID))
		},
	}

	cmdShowNext := &cobra.Command{
		Use:   "next",
		Short: "Find the next day with events and list them",
		RunE: func(cmd *cobra.Command, args []string) error {
			events, err := storage.Next(s, chatID)
			if err != nil {
				return err
			}
			fmt.Println(events)
			return nil
		},
	}

	cmdServe := &cobra.Command{
		Use:   "bot",
		Short: "Start Telegram bot to serve events and listen for updates",
		RunE: func(cmd *cobra.Command, args []string) error {
			return bot.ListenAndServe(token, s)
		},
	}

	cmd.AddCommand(cmdShow, cmdServe)
	cmd.PersistentFlags().StringVar(&eventsPath, "events", "events.yaml", "YAML file defining events")

	cmdShow.AddCommand(cmdShowAll, cmdShowNext)
	cmdShow.PersistentFlags().Int64Var(&chatID, "chat-id", 0, "Telegram chat ID")
	if err := cmdShow.MarkPersistentFlagRequired("chat-id"); err != nil {
		log.Fatalln("FATAL", err)
	}

	cmdServe.Flags().StringVar(&token, "token", "", "Telegram bot token")
	if err := cmdServe.MarkFlagRequired("token"); err != nil {
		log.Fatalln("FATAL", err)
	}

	_ = cmd.Execute()
}
