package main

import (
	"flag"
	"fmt"
	"github.com/Chef1991/gbot/bot"
	_ "github.com/Chef1991/gbot/plugins"
	"os"
)

var connectionString string
var teamName string

func init() {
	flag.StringVar(&connectionString, "conn", "", "Connection string example: http://bot@example.com:Password@gbot.chat:8080")
	flag.StringVar(&teamName, "team", "", "Name of team the bot will join")

}

func main() {
	flag.Parse()
	argError := false
	if len(connectionString) == 0 {
		fmt.Errorf("Connection string argument is required\n")
		argError = true
	}
	if len(teamName) == 0 {
		fmt.Errorf("Team name argument is required\n")
		argError = true
	}
	if argError {
		flag.PrintDefaults()
		os.Exit(1)
	}
	gbot, err := bot.NewBot(connectionString, teamName)
	if err != nil {
		fmt.Println("Error while creating bot")
		fmt.Println(err)
		panic(err)
	}
	gbot.Start()
	select {}

}
