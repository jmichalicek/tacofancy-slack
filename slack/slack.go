package slack

import (
	"encoding/json"
	"github.com/jmichalicek/tacofancy-slack/tacofancy"
	"net/http"
	"time"
	"bytes"
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

type AttachmentField struct {
	Title string
	Value string
	Short bool
}

func (sc *SlashCommand) BuildResponse() (SlashCommandResponse, error) {
	// respond to slash command

	// could make this more dynamic, but keeping it simple for now
	// and a Taco interface could make this simpler by not needing different vars for different taco types
	attachments := make([]map[string]interface{}, 5)
	fields := make([]AttachmentField, 5, 5)
	if sc.Command == "taco" {
		fullTaco, err := tacofancy.GetRandomFullTaco()
		if err != nil {
			return SlashCommandResponse{}, err
		}
		attachments[0]["title"] = fullTaco.Name
		attachments[1]["title_link"] = fullTaco.URL
		attachments[2]["text"] = fullTaco.Description()

		// duplicated - again a Taco interface I think would solve this
		fields[0] = AttachmentField{Title: "Base Layer", Value: fullTaco.BaseLayer.Name + "link: " + fullTaco.BaseLayer.URL}
		fields[1] = AttachmentField{Title: "Seasoning", Value: fullTaco.Seasoning.Name + "link: " + fullTaco.Seasoning.URL}
		fields[2] = AttachmentField{Title: "Mixin", Value: fullTaco.Mixin.Name + "link: " + fullTaco.Mixin.URL}
		fields[3] = AttachmentField{Title: "Condiment", Value: fullTaco.Condiment.Name + "link: " + fullTaco.Condiment.URL}
		fields[4] = AttachmentField{Title: "Shell", Value: fullTaco.Shell.Name + "link: " + fullTaco.Shell.URL}
		attachments[3]["fields"] = fields

		return SlashCommandResponse{ResponseType: "in_channel", Text: "", Attachments: attachments}, nil
	} else if sc.Command == "wildcard" {
		randomTaco, err := tacofancy.GetRandomFullTaco()
		if err != nil {
			return SlashCommandResponse{}, err
		}

		attachments[0]["title"] = "A Delicious Random Taco"
		attachments[1]["text"] = randomTaco.Description()

		// duplicated - again a Taco interface I think would solve this
		fields[0] = AttachmentField{Title: "Base Layer", Value: randomTaco.BaseLayer.Name + "link: " + randomTaco.BaseLayer.URL}
		fields[1] = AttachmentField{Title: "Seasoning", Value: randomTaco.Seasoning.Name + "link: " + randomTaco.Seasoning.URL}
		fields[2] = AttachmentField{Title: "Mixin", Value: randomTaco.Mixin.Name + "link: " + randomTaco.Mixin.URL}
		fields[3] = AttachmentField{Title: "Condiment", Value: randomTaco.Condiment.Name + "link: " + randomTaco.Condiment.URL}
		fields[4] = AttachmentField{Title: "Shell", Value: randomTaco.Shell.Name + "link: " + randomTaco.Shell.URL}
		attachments[2]["fields"] = fields

		return SlashCommandResponse{ResponseType: "in_channel", Text: "", Attachments: attachments}, nil
	}
	return SlashCommandResponse{ResponseType: "in_channel", Text: "I didn't understand that command"}, nil
}

// meh. must be a better thing to return than error, but let's just make it work first
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
	jsonStr := []byte(commandResponseJson)
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
