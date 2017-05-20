package bot

import (
	"net/http"
)

type SlashCommandHandler func(*Bot, SlashCommandRequest) *SlashCommandResponse

type SlashCommandRequest struct {
	ChannelId   string
	ChannelName string
	Command     string
	TeamDomain  string
	TeamId      string
	Text        string
	Token       string
	UserId      string
	UserName    string
}

// Message will show up in the channel as a normal message
const RESPONSE_TYPE_IN_CHANNEL = "in_channel"

// Message will appear temporarily and only display to the user who activated the command
const RESPONSE_TYPE_EPHEMERAL = "ephemeral"

type SlashCommandResponse struct {
	// The text to display in the message
	Text string `json:"text"`
	// The username to show as, if left blank will be set to the user who sent the message
	UserName string `json:"username"`
	// Determines how the Message will be displayed.  See RESPONSE_TYPE_IN_CHANNEL and RESPONSE_TYPE_EPHEMERAL
	ResponseType string `json:"response_type"`
	// Url to icon if you wish to override the default icon for the slash command
	IconUrl string `json:"icon_url"`
	// Webpage to redirect the user who sent the slash command to
	GotoLocation string `json:"goto_location"`
}

type SlashCommandConfig struct {
	Route       string
	ChannelName string
	Command     string
	Token       string
}

func NewSlashCommandResponse(req SlashCommandRequest) *SlashCommandResponse {
	resp := new(SlashCommandResponse)
	resp.ResponseType = RESPONSE_TYPE_IN_CHANNEL
	return resp
}

func ParseSlashCommandRequest(r *http.Request) *SlashCommandRequest {
	res := SlashCommandRequest{
		ChannelId:   r.FormValue("channel_id"),
		ChannelName: r.PostFormValue("channel_name"),
		Command:     r.PostFormValue("command"),
		TeamDomain:  r.PostFormValue("team_domain"),
		TeamId:      r.PostFormValue("team_id"),
		Text:        r.PostFormValue("text"),
		Token:       r.PostFormValue("token"),
		UserId:      r.PostFormValue("user_id"),
		UserName:    r.PostFormValue("user_name"),
	}
	if len(res.Command) == 0 {
		return nil
	}
	return &res
}
