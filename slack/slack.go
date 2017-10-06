package slack

import (
    "github.com/jmichalice/tacofancy-slack/tacofancy"
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
    Token string
    TeamId string
    TeamDomain string
    Enterprise string
    EnterpriseName string
    ChannelId string
    ChannelName string
    UserId string
    UserName string
    Command string
    Text string
    ResponseURL string
    TriggerId string
}

type AttachmentField struct {
    Title string
    Value string
    Short bool
}

func (sc *SlashCommand) Respond() (SlashCommandResponse, error) {
    // respond to slash command

    // could make this more dynamic, but keeping it simple for now
    // and a Taco interface could make this simpler by not needing different vars for different taco types
    attachments := make([]map[string]interface{})
    fields = make([]AttachmentField, 5, 5)
    if sc.Command == "taco" {
        fullTaco, err := tacofancy.GetRandomFullTaco()
        attachments["title"] = fullTaco.Name
        attachments["title_link"] = fullTaco.URL
        attachments["text"] = fullTaco.Description()

        // duplicated - again a Taco interface I think would solve this
        fields[0] = AttachmentField{Name: "Base Layer", Value: fullTaco.BaseLayer.Name + "link: " + fullTaco.BaseLayer.URL}
        fields[1] = AttachmentField{Name: "Seasoning", Value: fullTaco.Seasoning.Name + "link: " + fullTaco.Seasoning.URL}
        fields[2] = AttachmentField{Name: "Mixin", Value: fullTaco.Mixin.Name + "link: " + fullTaco.Mixin.URL}
        fields[3] = AttachmentField{Name: "Condiment", Value: fullTaco.Condiment.Name + "link: " + fullTaco.Condiment.URL}
        fields[4] = AttachmentField{Name: "Shell", Value: fullTaco.Shell.Name + "link: " + fullTaco.Shell.URL}
        attachments["fields"] = fields

        return SlashCommandResponse{
            ResponseType: "in_channel", Text: "",
            Attachments: attachments}, nil
    } else if sc.Command == "wildcard" {
        randomTaco, err := tacofancy.GetRandomFullTaco()
        attachments["title"] = "A Delicious Random Taco"
        attachments["text"] = randomTaco.Description()

        // duplicated - again a Taco interface I think would solve this
        fields[0] = AttachmentField{Name: "Base Layer", Value: fullTaco.BaseLayer.Name + "link: " + fullTaco.BaseLayer.URL}
        fields[1] = AttachmentField{Name: "Seasoning", Value: fullTaco.Seasoning.Name + "link: " + fullTaco.Seasoning.URL}
        fields[2] = AttachmentField{Name: "Mixin", Value: fullTaco.Mixin.Name + "link: " + fullTaco.Mixin.URL}
        fields[3] = AttachmentField{Name: "Condiment", Value: fullTaco.Condiment.Name + "link: " + fullTaco.Condiment.URL}
        fields[4] = AttachmentField{Name: "Shell", Value: fullTaco.Shell.Name + "link: " + fullTaco.Shell.URL}
        attachments["fields"] = fields

        return SlashCommandResponse{
            ResponseType: "in_channel", Text: ""
            Attachments: attachments}, nil
    } else {
        return SlashCommandResponse{ResponseType: "in_channel", Text: "I didn't understand that command"}, nil
    }
}

type SlashCommandResponse struct {
    ResponseType string `json:"response_type"`
    Text string `json:"text"`
    // attachments could get more complex
    // https://api.slack.com/docs/message-attachments
    // using interface as the value data type here because it could be strings or a list
    Attachments []map[string]interface{} `json:"attachments"`
}

func SendDelayedResponse(sc SlashCommand) error {
    var httpClient = &http.Client{
    	Timeout: time.Second * 10,
    }
    url := sc.ResponsseURL
    commandResponse := sc.Respond()
    commandResponseJson, err := json.Marshal(scr, "", "  ")
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

    return
}
