package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var version string

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
		Version:           version,
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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
		Short: "List events",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(events)
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
				bot.send(events.today()...)
			}
		},
	}

	cmd.AddCommand(cmdList, cmdServe)
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
