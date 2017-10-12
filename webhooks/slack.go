package webhooks

// should I make a webhooks/slack package instead?

import (
	"github.com/jmichalicek/tacofancy-slack/slack"
	"net/http"
)

func SlashCommandHandler(w http.ResponseWriter, r *http.Request) {
	// does not come in as json, but response should be json
	command, text, token := r.FormValue("command"), r.FormValue("text"), r.FormValue("token")
	responseURL, channelId, userId := r.FormValue("response_url"), r.FormValue("channel_id"), r.FormValue("user_id")
	slashCommand := slack.SlashCommand{
		Command: command, Text: text, Token: token, ResponseURL: responseURL, ChannelId: channelId, UserId: userId}

	if slack.VerifyToken(slashCommand.Token) {
		// could move some of this to the SlashCommand
		// or should this use sc.RespondAsync()?
		scr, err := sc.BuildResponse()
		if err == nil {
			// //TODO: how to test this?
			go slack.SendDelayedResponse(sc.URL(), scr)
			w.WriteHeader(200)
			return
		}
	}
	w.WriteHeader(400)
}
