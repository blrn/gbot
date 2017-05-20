package bot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type RouteHandler func(*Bot, http.ResponseWriter, *http.Request)

type Route struct {
	path      *string
	method    string
	regexPath *regexp.Regexp
	handler   RouteHandler
}
type slashCommand struct {
	config  SlashCommandConfig
	handler SlashCommandHandler
}
type HttpRouter struct {
	routes        []Route
	slashCommands []slashCommand
	gBot          *Bot
}

func (router *HttpRouter) Get(path string, handler RouteHandler) {
	router.newRoute(http.MethodGet, path, handler)
}

func (router *HttpRouter) Post(path string, handler RouteHandler) {
	router.newRoute(http.MethodPost, path, handler)
}

func (router *HttpRouter) GetRegex(path *regexp.Regexp, handler RouteHandler) {
	router.newRegexRoute(http.MethodGet, path, handler)
}
func (router *HttpRouter) PostRegex(path *regexp.Regexp, handler RouteHandler) {
	router.newRegexRoute(http.MethodPost, path, handler)
}
func (router *HttpRouter) newRegexRoute(method string, pathRegex *regexp.Regexp, handler RouteHandler) {
	r := Route{method: method, regexPath: pathRegex, handler: handler}
	router.addRoute(r)
}

func (router *HttpRouter) newRoute(method string, path string, handler RouteHandler) {
	r := Route{method: method, path: &path, handler: handler}
	router.addRoute(r)
}
func (router *HttpRouter) addRoute(route Route) {
	if router.routes == nil {
		router.routes = make([]Route, 0, 10)
	}
	router.routes = append(router.routes, route)
}
func (router *HttpRouter) addSlashCommand(cmdConfig SlashCommandConfig, handler SlashCommandHandler) {
	if router.slashCommands == nil {
		router.slashCommands = make([]slashCommand, 0, 10)
	}
	router.slashCommands = append(router.slashCommands, slashCommand{config: cmdConfig, handler: handler})
}

func (router HttpRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.ToLower(r.URL.Path)
	handled := false
	if router.slashCommands != nil {
		if sCmd := ParseSlashCommandRequest(r); sCmd != nil {
			for _, cmd := range router.slashCommands {
				if strings.ToLower(cmd.config.Command) != strings.ToLower(sCmd.Command) {
					continue
				}
				if len(cmd.config.ChannelName) > 0 && strings.ToLower(cmd.config.ChannelName) != strings.ToLower(sCmd.ChannelName) {
					continue
				}
				if len(cmd.config.Token) > 0 && strings.ToLower(cmd.config.Token) != strings.ToLower(sCmd.Token) {
					continue
				}
				// looks like it matches
				if resp := cmd.handler(router.gBot, *sCmd); resp != nil {
					if respMarshalled, err := json.Marshal(resp); err != nil {
						fmt.Println("ERROR marshalling response", err)
					} else {
						w.Write(respMarshalled)
					}
				}
				handled = true
				break
			}
		}
	}
	if router.routes != nil && !handled {
		for _, route := range router.routes {
			if route.method == r.Method {
				if route.path != nil && strings.Contains(path, strings.ToLower(*route.path)) {
					route.handler(router.gBot, w, r)
					handled = true
					break
				} else if route.regexPath != nil && route.regexPath.MatchString(path) {
					route.handler(router.gBot, w, r)
					handled = true
					break
				}
			}
		}
	}
	if !handled {
		fmt.Fprintf(w, "Nope\n")
	}
}
