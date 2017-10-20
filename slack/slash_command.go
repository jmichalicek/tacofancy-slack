package slack;

import (
	"github.com/jmichalicek/tacofancy-slack/tacofancy"
	"strings"
)

// Returns a new SlashCommand struct
func NewSlashCommand(token, command, text, responseURL, channelID, channelName, userID, userName, teamID, teamDomain,
	enterprise, enterpriseName, triggerID string, tacofancyClient tacofancy.Client) SlashCommand {

	return SlashCommand{token: token, teamID: teamID, teamDomain: teamDomain, enterprise: enterprise,
		enterpriseName: enterpriseName, channelID: channelID, userName: userName, command: command, text: text,
		responseURL: responseURL, triggerID: triggerID, tacofancyClient: tacofancyClient}

}

// Example slash command from docs
// token=gIkuvaNzQIHg97ATvDxqgjtO
// team_id=T0001
// team_domain=example
// enterprise_id=E0001
// enterprise_name=Globular%20Construct%20Inc
// channel_id=C2147483705
// channel_name=test
// user_id=U2147483697
// user_name=Steve
// command=/weather
// text=94070
// response_url=https://hooks.slack.com/commands/1234/5678
// trigger_id=13345224609.738474920.8088930838d88f008e0
// https://api.slack.com/slash-commands
type SlashCommand struct {
	token          string
	teamID         string
	teamDomain     string
	enterprise     string
	enterpriseName string
	channelID      string
	channelName    string
	userID         string
	userName       string
	command        string
	text           string
	responseURL    string
	triggerID      string

	// I do not really feel like this belongs here, but for the sake of easy testability
	// it will live here for now.  In a more complex app perhaps there would be a
	// TacoFancySlashCommand where this would feel more appropriate
	// or perhaps BuildResponse() should not be a SlashCommand method and I should pass
	// the SlashCommand and client into NewTacoRecipeResponse() etc then a factory method
	// which is basiclaly what BuildResponse() is.  The problem here is that I feel like
	// no matter what, at some point, I am configuring a tacofancy.Client and passing it inot
	// a function, not because I know the actual code needs it, but because the code "might" need it
	// and test cases need some way to pass it in for those.
	tacofancyClient tacofancy.Client
}

func (sc SlashCommand) Token() string                     { return sc.token }
func (sc SlashCommand) TeamID() string                    { return sc.teamID }
func (sc SlashCommand) TeamDomain() string                { return sc.teamDomain }
func (sc SlashCommand) Enterprise() string                { return sc.enterprise }
func (sc SlashCommand) EnterpriseName() string            { return sc.enterpriseName }
func (sc SlashCommand) ChannelID() string                 { return sc.channelID }
func (sc SlashCommand) ChannelName() string               { return sc.channelName }
func (sc SlashCommand) UserID() string                    { return sc.userID }
func (sc SlashCommand) UserName() string                  { return sc.userName }
func (sc SlashCommand) Command() string                   { return sc.command }
func (sc SlashCommand) Text() string                      { return sc.text }
func (sc SlashCommand) ResponseURL() string               { return sc.responseURL }
func (sc SlashCommand) TriggerID() string                 { return sc.triggerID }
func (sc SlashCommand) TacofancyClient() tacofancy.Client { return sc.tacofancyClient }

// Builds a SlashCommandResponse to use for responding to a Slack SlashCommand
// TODO:  make a responder interface and pass that in?  then t.Respond() could
// be called?  Not sure since that moves the slack specific logic out of here
// but could almost certainly be far more useful for a more general slack framework
func (sc SlashCommand) BuildResponse() (SlashCommandResponse, error) {
	// respond to slash command
	// TODO: Have a response factory registry?
	parts := strings.Split(sc.Text(), " ")
	commandType := "recipe"
	if len(parts) >= 1 && parts[0] != "" {
		commandType = strings.ToLower(parts[0])
	}

	if commandType == "recipe" {
		return NewTacoRecipeResponse(sc.TacofancyClient())
	} else if commandType == "loco" {
		return NewTacoLocoResponse(sc.TacofancyClient())
	} else if commandType == "grande" {
		return NewTacoGrandeResponse()
	}
	return NewSlashCommandResponse("in_channel", "I didn't understand that command", nil), nil
}

// Builds a SlashCommandResponse and calls SendDelayedResponse() and returns the SendDelayedResponse()'s
// return value in the channel.
func (sc SlashCommand) RespondAsync(ch chan error) {
	defer close(ch)
	scr, _ := sc.BuildResponse()
	ch <- SendDelayedResponse(sc.ResponseURL(), scr)
}
