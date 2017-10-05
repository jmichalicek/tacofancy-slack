package slack

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

func (sc *SlashCommand) Respond() (SlashCommandResponse, error) {
    // respond to slash command

    // could make this more dynamic, but keeping it simple for now
    if sc.Command == "taco" {
        return SlashCommandResponse{ResponseType: "in_channel", Text: "I will eventually link to a taco recipe"}, nil
    } else {
        return SlashCommandResponse{ResponseType: "in_channel", Text: "I didn't understand that command"}, nil
    }
}

type SlashCommandResponse struct {
    ResponseType string `json:"response_type"`
    Text string `json:"text"`
    // attachments could get more complex
    // https://api.slack.com/docs/message-attachments
    Attachments []map[string]string `json:"attachments"`
}
