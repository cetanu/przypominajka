package main

import (
	"fmt"
	"log"
	"time"

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
		events     events
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
			e, err := readEvents(eventsPath)
			if err != nil {
				return err
			}
			events = e
			return nil
		},
	}

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List all events",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(events)
		},
	}

	cmdNext := &cobra.Command{
		Use:   "next",
		Short: "Find the next day with events and list them",
		Run: func(cmd *cobra.Command, args []string) {
			month, day, next := events.next()
			if len(next) == 0 {
				fmt.Println("No events found")
				return
			}
			for _, event := range next {
				fmt.Printf(formatListLine, day, month, event)
			}
		},
	}

	cmdServe := &cobra.Command{
		Use:   "serve",
		Short: "Start Telegram bot",
		Run: func(cmd *cobra.Command, args []string) {
			bot, err := newBot(token, chatID)
			if err != nil {
				log.Fatalln("FATAL", err)
			}
			go bot.listen()

			for t := range time.Tick(time.Hour) {
				if t.Round(time.Hour).Hour() != 9 { // run once a day between 8:30 and 9:29
					continue
				}
				for _, e := range events.today() {
					if err := bot.send(e); err != nil {
						log.Println("ERROR", err)
					}
				}
			}
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
