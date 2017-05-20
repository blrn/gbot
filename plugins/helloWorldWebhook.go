package plugins

import (
	"fmt"
	"github.com/blrn/gbot/bot"
	"net/http"
)

func init() {
	bot.Register("githubDescription", func(gbot *bot.Bot) {
		gbot.Router.Get("/hello/world", OnHelloWorld)
		gbot.Router.Post("/hello/world", OnHelloWorld)
	})
}

func OnHelloWorld(gBot *bot.Bot, w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method == http.MethodGet)
	bot.ParseSlashCommandRequest(r)
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		name = "World"
	}
	p := gBot.NewPost()
	fmt.Println(p.ChannelId)
	p.Message = fmt.Sprintf("Hello, %s!!!", name)
	gBot.SendMessage(p)
	fmt.Fprintf(w, "Ok")
}
