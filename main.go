package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var version string

const description = `przypominajka - a Telegram bot for sending event reminders

Description:
  przypominajka reads a YAML file with events and sends reminders about them.
  The reminders are sent out on the day of the event between 08:30 and 09:29
  system time (exact time depends on serve command startup time).

  Reminders are written in Polish.

Example events.yaml:
  january:
    5:
      - name: "John"
        type: "birthday"
      - name: "Jane"
        surname: "Doe"
        type: "nameday"
  april:
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
		year       year
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
			if use := cmd.Parent().Use; use == "completion" || use == "help" {
				return nil
			}
			y, err := readYear(eventsPath)
			if err != nil {
				return err
			}
			year = y
			return nil
		},
	}

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List all events",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(year)
		},
	}

	cmdNext := &cobra.Command{
		Use:   "next",
		Short: "Find the next day with events and list them",
		Run: func(cmd *cobra.Command, args []string) {
			month, day, next := year.next()
			if len(next) == 0 {
				fmt.Println("No events found")
				return
			}
			fmt.Println(next.Format(month, day))
		},
	}

	cmdServe := &cobra.Command{
		Use:   "bot",
		Short: "Start Telegram bot to serve events and listen for updates",
		Run: func(cmd *cobra.Command, args []string) {
			bot, err := newBot(token, chatID, year)
			if err != nil {
				log.Fatalln("FATAL", err)
			}
			go bot.listen()
			bot.serve()
		},
	}

	cmd.AddCommand(cmdList, cmdNext, cmdServe)
	cmd.PersistentFlags().StringVar(&eventsPath, "events", "events.yaml", "YAML file defining events")
	cmdServe.Flags().StringVar(&token, "token", "", "Telegram bot token")
	if err := cmdServe.MarkFlagRequired("token"); err != nil {
		log.Fatalln("FATAL", err)
	}
	cmdServe.Flags().Int64Var(&chatID, "chat-id", 0, "Telegram chat ID")
	if err := cmdServe.MarkFlagRequired("chat-id"); err != nil {
		log.Fatalln("FATAL", err)
	}

	_ = cmd.Execute()
}
