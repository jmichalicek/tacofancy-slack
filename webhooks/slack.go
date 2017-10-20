package webhooks

// should I make a webhooks/slack package instead?

import (
	"github.com/jmichalicek/tacofancy-slack/slack"
	"github.com/jmichalicek/tacofancy-slack/tacofancy"
	"net/http"
)

func SlashCommandHandler(w http.ResponseWriter, r *http.Request) {
	// does not come in as json, but response should be json
	command, text, token := r.FormValue("command"), r.FormValue("text"), r.FormValue("token")
	responseURL, channelID, userID := r.FormValue("response_url"), r.FormValue("channel_id"), r.FormValue("user_id")

	client := tacofancy.NewClient("", nil)
	slashCommand := slack.NewSlashCommand(token, command, text, responseURL, channelID, "", userID, "", "", "", "", "",
		"", client)

	// for a more generic implementation, perhaps the token could be used to look up
	// the slack app and an appropriate SlashCommand
	if slack.VerifyToken(slashCommand.Token()) {
		// could move some of this to the SlashCommand
		// or should this use sc.RespondAsync()?
		scr, err := slashCommand.BuildResponse()
		if err == nil {
			// //TODO: how to test this?
			go slack.SendDelayedResponse(slashCommand.ResponseURL(), scr)
			w.WriteHeader(200)
			return
		}
	}
	w.WriteHeader(400)
}
