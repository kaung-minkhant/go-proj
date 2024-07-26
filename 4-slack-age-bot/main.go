package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/slack-io/slacker"
)

// func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
// 	for event := range analyticsChannel {
// 		fmt.Println("Command Events")
// 		fmt.Println(event.Timestamp)
// 		fmt.Println(event.Command)
// 		fmt.Println(event.Parameters)
// 		fmt.Println(event.Event)
// 		fmt.Println()
// 	}
// }

func main() {
	os.Setenv("SLACK_BOT_TOKEN", "")
	os.Setenv("SLACK_APP_TOKEN", "")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	bot.AddCommand(&slacker.CommandDefinition{
		Command: "ping",
		Handler: func(cc *slacker.CommandContext) {
			cc.Response().Reply("pong")
		},
	})

	bot.AddCommand(&slacker.CommandDefinition{
		Command:     "my yob is {year}",
		Description: "year of birth calculator",
		Examples:    []string{"my yob is 2024"},
		Handler: func(cc *slacker.CommandContext) {
			year := cc.Request().Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				fmt.Println("error: ", err)
				return
			}
			age := time.Now().Year() - yob
			response := fmt.Sprintf("Your age is %d", age)
			cc.Response().Reply(response)
		},
	})

	// bot.Command("my yob is <year>", &slacker.CommandDefinition{
	// 	Description: "yob calculator",
	// 	Examples:    []string{"my yob is 2024"},
	// 	Handler: func(bc slacker.BotContext, r slacker.Request, w slacker.ResponseWriter) {
	// 		year := r.Param("year")
	// 		yob, err := strconv.Atoi(year)
	// 		if err != nil {
	// 			fmt.Println("error")
	// 		}
	// 		age := time.Now().Year() - yob
	// 		response := fmt.Sprintf("Your age is %d", age)
	// 		w.Reply(response)
	// 	},
	// })

	// go printCommandEvents(bot.CommandEvents())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
