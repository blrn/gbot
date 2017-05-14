package bot

type RegisterRequest struct {
	Name      string
	ReadyFunc func(*Bot)
}

var registeredBots []RegisterRequest

func init() {
	registeredBots = make([]RegisterRequest, 0, 10)
}

func Register(name string, ready func(bot *Bot)) {
	registeredBots = append(registeredBots, RegisterRequest{Name: name, ReadyFunc: ready})
}

func registerPlugins(bot *Bot) {
	for _, registeredBot := range registeredBots {
		registeredBot.ReadyFunc(bot)
	}
}
