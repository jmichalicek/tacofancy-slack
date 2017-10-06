package webhooks

// should I make a webhooks/slack package instead?

import (
	// "encoding/json"
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
		go sendAsynResponse(slashCommand)
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}

}

func sendAsynResponse(sc slack.SlashCommand) {
	// could be interesting to do the slashCommand part in a channel?
	// unnecessary, just for learning/fiddling about.
	// TODO: handle errors
	scr, _ := sc.BuildResponse()
	slack.SendDelayedResponse(sc, scr)
}
