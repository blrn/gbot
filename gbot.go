package main

import (
	"fmt"
	"github.com/Chef1991/gbot/bot"
	_ "github.com/Chef1991/gbot/plugins"
)

func main() {
	gbot, err := bot.NewBot("http://gbot@tcook.me:Test123$@192.168.1.12", "testTeam")
	if err != nil {
		fmt.Println("Error while creating bot")
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(gbot)
	gbot.Start()
	select {}

}
