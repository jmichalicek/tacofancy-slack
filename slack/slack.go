package slack

import (
	"bytes"
	"encoding/json"
	"github.com/jmichalicek/tacofancy-slack/tacofancy"
	"net/http"
	"os"
	"strings"
	"time"
)

// Example slash command from docs
//     token=gIkuvaNzQIHg97ATvDxqgjtO
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
	Token          string
	TeamId         string
	TeamDomain     string
	Enterprise     string
	EnterpriseName string
	ChannelId      string
	ChannelName    string
	UserId         string
	UserName       string
	Command        string
	Text           string
	ResponseURL    string
	TriggerId      string
}

// A slashcommand response attachment
type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool `json:"short"`
	TitleLink string `json:"title_link"`
}

// Verifies the slack command's token by matching it up to
// the environment variable TACOFANCY_VERIFICATION_TOKEN
// Returns true if the tokens match, otherwise returns false
func VerifyToken(token string) bool {
	return token == os.Getenv("TACOFANCY_VERIFICATION_TOKEN")
}

// Builds a SlashCommandResponse to use for responding to a Slack SlashCommand
func (sc *SlashCommand) BuildResponse() (SlashCommandResponse, error) {
	// respond to slash command

	// could make this more dynamic, but keeping it simple for now
	// and a Taco interface could make this simpler by not needing different vars for different taco types
	attachments := make([]map[string]interface{}, 1)
	attachments[0] = make(map[string]interface{})
	fields := make([]AttachmentField, 5, 5)

	// TODO: split these out into separate handlers?
	parts := strings.Split(sc.Text, " ")
	commandType := "recipe"
	if len(parts) >= 1 && parts[0] != "" {
		commandType = strings.ToLower(parts[0])
	}

	if commandType == "recipe" {
		fullTaco, err := tacofancy.GetRandomFullTaco()
		if err != nil {
			return SlashCommandResponse{}, err
		}
		// A bunch of duplicated stuff which it seems like using Taco interface would solve
		// but since Taco interface cannot access the struct properties, here we are... unless the two objects
		// just become one with some unused parts or I go with a whole bunch of duplicated getters/setters
		attachments[0]["title"] = fullTaco.Name
		attachments[0]["title_link"] = fullTaco.URL
		attachments[0]["text"] = fullTaco.Description()
		if fullTaco.BaseLayer.Name != "" {
			fields[0] = AttachmentField{Title: "Base Layer", Value: "<"+fullTaco.BaseLayer.URL+"|"+fullTaco.BaseLayer.Name+">", Short: true}
		}
		if fullTaco.Seasoning.Name != "" {
			fields[1] = AttachmentField{Title: "Seasoning: ", Value: "<"+fullTaco.Seasoning.URL+"|"+fullTaco.Seasoning.Name+">", Short: true}
		}
		if fullTaco.Mixin.Name != "" {
			fields[2] = AttachmentField{Title: "Mixin: ", Value: "<"+fullTaco.Mixin.URL+"|"+fullTaco.Mixin.Name+">", Short: true}
		}
		if fullTaco.Condiment.Name != "" {
			fields[3] = AttachmentField{Title: "Condiment: ", Value: "<"+fullTaco.Condiment.URL+"|"+fullTaco.Condiment.Name+">", Short: true}
		}
		if fullTaco.Shell.Name != "" {
			fields[4] = AttachmentField{Title: "Shell: ", Value: "<"+fullTaco.Shell.URL+"|"+fullTaco.Shell.Name+">", Short: true}
		}
		attachments[0]["fields"] = fields

		return SlashCommandResponse{ResponseType: "in_channel", Text: "", Attachments: attachments}, nil
	} else if commandType == "loco" {
		randomTaco, err := tacofancy.GetRandomTacoParts()
		if err != nil {
			return SlashCommandResponse{}, err
		}
		attachments[0]["title"] = "A Delicious Random Taco"
		attachments[0]["text"] = randomTaco.Description()
		if randomTaco.BaseLayer.Name != "" {
			fields[0] = AttachmentField{Title: "Base Layer", Value: "<"+randomTaco.BaseLayer.URL+"|"+randomTaco.BaseLayer.Name+">", Short: true}
		}
		if randomTaco.Seasoning.Name != "" {
			fields[1] = AttachmentField{Title: "Seasoning: ", Value: "<"+randomTaco.Seasoning.URL+"|"+randomTaco.Seasoning.Name+">", Short: true}
		}
		if randomTaco.Mixin.Name != "" {
			fields[2] = AttachmentField{Title: "Mixin: ", Value: "<"+randomTaco.Mixin.URL+"|"+randomTaco.Mixin.Name+">", Short: true}
		}
		if randomTaco.Condiment.Name != "" {
			fields[3] = AttachmentField{Title: "Condiment: ", Value: "<"+randomTaco.Condiment.URL+"|"+randomTaco.Condiment.Name+">", Short: true}
		}
		if randomTaco.Shell.Name != "" {
			fields[4] = AttachmentField{Title: "Shell: ", Value: "<"+randomTaco.Shell.URL+"|"+randomTaco.Shell.Name+">", Short: true}
		}
		attachments[0]["fields"] = fields

		return SlashCommandResponse{ResponseType: "in_channel", Text: "", Attachments: attachments}, nil
	}
	return SlashCommandResponse{ResponseType: "in_channel", Text: "I didn't understand that command"}, nil
}

// Builds a SlashCommandResponse and calls SendDelayedResponse() and returns the SendDelayedResponse()'s
// return value in the channel.
func (sc *SlashCommand) RespondAsync(ch chan error) {
	defer close(ch)
	scr, _ := sc.BuildResponse()
	ch <- SendDelayedResponse(*sc, scr)
}

type SlashCommandResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
	// attachments could get more complex
	// https://api.slack.com/docs/message-attachments
	// using interface as the value data type here because it could be strings or a list
	Attachments []map[string]interface{} `json:"attachments"`
}

// Sends the provided SlashCommandResponse to the provided SlashCommand
// by Marshalling the SlashCommandResponse and making an HTTP POST to the SlashCommand.ResponseURL
func SendDelayedResponse(sc SlashCommand, scr SlashCommandResponse) error {
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
	url := sc.ResponseURL
	//commandResponse, _ := sc.BuildResponse()
	commandResponseJson, err := json.Marshal(scr)
	if err != nil {
		return err
	}

	// jsonStr and bytes.NewBuffer() from https://stackoverflow.com/a/24455606
	// TODO: See if there is a better way
	// not sure I need this.. isn't it already bytes when it comes out of Marshal?
	// jsonStr := []byte(commandResponseJson)
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(commandResponseJson))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
