package bot

import (
	"errors"
	"fmt"
	"github.com/mattermost/platform/model"
	"net/url"
	"regexp"
	"strings"
)

func init() {
	fmt.Println("Hello from bot!")
}

type Plugin struct {
	isRegex       bool
	matchStr      string
	onMessage     func(post *model.Post, bot *Bot)
	compiledRegex *regexp.Regexp
}

type Bot struct {
	client           *model.Client
	webSocketClient  *model.WebSocketClient
	user             *model.User
	initialLoad      *model.InitialLoad
	defaultChannelId string
	serverUrl        string
	websocketUrl     string
	userInfo         *url.Userinfo
	team             *model.Team
	plugins          []Plugin
	stopChan         chan bool
}

// Create a new Bot Object
//
// connection string is a url string containing the username, password, server address and path to connect to
// the mattermost server
//
// Returns a pointer to a Bot object and an error
func NewBot(connectionString string, teamName string) (*Bot, error) {
	serverUrl, err := url.Parse(connectionString)
	if err != nil {
		return nil, err
	}
	scheme := serverUrl.Scheme
	username := serverUrl.User.Username()
	password, hasPass := serverUrl.User.Password()
	host := serverUrl.Host
	opaque := serverUrl.Opaque
	fragment := serverUrl.Fragment
	path := serverUrl.Path
	port := serverUrl.Port()
	fmt.Printf("Scheme: %s\nusername: %s\nhasPass: %t\nPassword: %s\nHost: %s\nOpaque: %s\nfragment: %s\npath: %s\nport: %s\n", scheme, username, hasPass, password, host, opaque, fragment, path, port)
	bot := new(Bot)
	bot.serverUrl = fmt.Sprintf("%s://%s%s", serverUrl.Scheme, serverUrl.Host, serverUrl.Path)
	bot.websocketUrl = fmt.Sprintf("ws://%s%s", serverUrl.Host, serverUrl.Path)
	bot.userInfo = serverUrl.User
	if err = bot.createClient(); err != nil {
		return bot, err
	}
	if err = bot.login(); err != nil {
		return bot, err
	}
	if initLoad, err := bot.client.GetInitialLoad(); err != nil {
		return bot, err
	} else {
		bot.initialLoad = initLoad.Data.(*model.InitialLoad)
	}
	if err = bot.SetTeam(teamName); err != nil {
		return bot, err
	}
	if err = bot.createWsClient(); err != nil {
		return bot, err
	}
	registerPlugins(bot)
	return bot, nil
}
func (bot *Bot) SetTeam(teamName string) error {
	for _, team := range bot.initialLoad.Teams {
		if strings.ToLower(teamName) == strings.ToLower(team.Name) {
			bot.team = team
			bot.client.SetTeamId(bot.team.Id)
			break
		}
	}
	if bot.team == nil {
		return errors.New("Could not find team")
	} else {
		return nil
	}
}
func (bot *Bot) createWsClient() error {
	if webSocketClient, err := model.NewWebSocketClient(bot.websocketUrl, bot.client.AuthToken); err != nil {
		return err
	} else {
		bot.webSocketClient = webSocketClient
		return nil
	}
}
func (bot *Bot) login() error {
	pass, hassPass := bot.userInfo.Password()
	if !hassPass {
		return errors.New("No password for bot account")
	}
	if loginResult, err := bot.client.Login(bot.userInfo.Username(), pass); err != nil {
		return err
	} else {
		bot.user = loginResult.Data.(*model.User)
		return nil
	}
}
func (bot *Bot) createClient() error {
	bot.client = model.NewClient(bot.serverUrl)
	if _, err := bot.client.GetPing(); err != nil {
		return err
	}
	return nil
}
func (bot *Bot) Match(callback func(*model.Post, *Bot), matchStrs ...string) {
	if bot.plugins == nil {
		bot.plugins = make([]Plugin, 0, 10)
	}
	for _, match := range matchStrs {
		plug := Plugin{matchStr: match, onMessage: callback, isRegex: false}
		bot.plugins = append(bot.plugins, plug)
	}
}
func (bot *Bot) MatchRegex(callback func(*model.Post, *Bot), matchRegex ...string) {
	if bot.plugins == nil {
		bot.plugins = make([]Plugin, 0, 10)
	}
	for _, reg := range matchRegex {
		compiled := regexp.MustCompile(reg)
		plug := Plugin{isRegex: true, matchStr: reg, compiledRegex: compiled, onMessage: callback}
		bot.plugins = append(bot.plugins, plug)
	}
}

func (bot *Bot) SendMessage(post *model.Post) error {
	if _, err := bot.client.CreatePost(post); err != nil {
		return err
	} else {
		return nil
	}
}

func (bot *Bot) Start() error {
	if bot.webSocketClient == nil {
		if err := bot.createWsClient(); err != nil {
			return err
		}
	}
	bot.webSocketClient.Listen()
	go func() {
		for {
			select {
			case event := <-bot.webSocketClient.EventChannel:
				bot.handleEvent(event)
			case <-bot.stopChan:
				break
			}
		}
	}()
	return nil
}
func (bot *Bot) handleEvent(event *model.WebSocketEvent) {
	fmt.Printf("handleEvent(%+v)\n", *event)
	if event == nil {
		return
	}
	if event.Event == model.WEBSOCKET_EVENT_POSTED {
		post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
		bot.handleNewPost(post)
	}

}
func (bot *Bot) handleNewPost(post *model.Post) {
	if post.UserId != bot.user.Id {
		fmt.Printf("handleNewPost(%+v)\n", *post)
		if bot.plugins != nil {
			for _, plug := range bot.plugins {
				if plug.isRegex {
					if plug.compiledRegex.MatchString(post.Message) {
						callback := plug.onMessage
						callback(post, bot)
					}
				}
			}
		}
	}
}
