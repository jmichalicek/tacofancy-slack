package webhooks
// should I make a webhooks/slack package instead?

import (
	"encoding/json"
	"github.com/jmichalicek/tacofancy-slack/slack"
	"net/http"
)

func SlashCommandHandler(w http.ResponseWriter, r *http.Request) {
	// does not come in as json, but response should be json
	command, text, token := r.FormValue("command"), r.FormValue("text"), r.FormValue("token")
	slashCommand := &slack.SlashCommand{Command: command, Text: text, Token: token}

	// could move some of this to the SlashCommand
	responseData, _ := slashCommand.Respond()
	responseJSON, _ := json.Marshal(responseData)

	// TODO: Parse the incoming json and do stuff
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(responseJSON)
}
