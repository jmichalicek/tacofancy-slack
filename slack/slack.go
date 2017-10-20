package slack

import (
	"bytes"
	"encoding/json"
	"github.com/jmichalicek/tacofancy-slack/tacofancy"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// converts https://raw.githubusercontent.com/sinker/tacofancy/master/condiments/baja_white_sauce.md
// to https://github.com/sinker/tacofancy/blob/master/condiments/baja_white_sauce.md
func githubRawUrlToRepo(url string) string {
	url = strings.Replace(url, "raw.github", "github", -1)
	url = strings.Replace(url, "master", "blob/master", -1)
	url = strings.Replace(url, "master//", "master/", -1)
	return url
}

var tacoQuotes []string = []string{
	"Yo quiero chimichangas y chile colorado\nYo tengo el dinero para un steak picado\nLas flautas y tamales, siempre muy bueno\nY el chile relleno",
	"You see, I just gotta have a tostada, carne asada\nThat's right, I want the whole enchilada\nMy only addiction has to do with a flour tortilla\nI need a quesadilla",
	"I love to stuff my face with tacos al carbón\nWith my friends, or when I'm all alone\nYo tengo mucho hambre y ahora lo quiero\nUn burrito ranchero",
	"So give me something spicy and hot, now\nBreak out the menu, what you got, now?\nOh, would you tell the waiter I'd like to have sour cream on the side\nYou better make sure the beans are refried",
	"Well, there's not a taco big enough for a man like me\nThat's why I order two or three\nLet me give you a tip, just try a nacho chip\nIt's really good with bean dip",
	"I eat uno, dos, tres, quatro burritos\nPretty soon I can't fit in my Speedos\nWell, I hope they feed us lots of chicken fajitas\nAnd a pitcher of margaritas",
	"Well, the combination plates all come with beans and rice\nThe taquitos here are very nice\nNow I'm down on my knees, we need some extra tomatoes and cheese\nAnd could you make that separate checks, please?",
	"Well, the food is coming, I can hardly wait\nNow watch your fingers, careful hot plate!\nWhat you think you're doing with my chile con queso?\nWell, if you want some, just say so",
	"Oh boy, pico de gallo\nThey sure don't make it like this in Ohio\nNo gracias, yo quiero jalapeños, nada más\nYou can toss away the hot sauce",
	"Donde estan los nachos? Holy frijole!\nYou better get me a bowl of guacamole\nY Usted, Eugene? Why's your face turning green?\nDon't you like pinto beans?",
	"You want some more cinnamon crispas?\nIf you don't, hasta la vista\nJust take the rest home in a doggie bag if you wanna\nYou can finish it mañana",
	"Well, it's been a pleasure, I can't eat no more\nSeñor, la cuenta, por favor\nIf you ain't ever tried real Mexican cooking, well, you oughta\nJust don't drink the water",
}

func getQuote() string {
	random := func(min, max int) int {
		rand.Seed(time.Now().Unix())
		return rand.Intn(max-min) + min
	}

	r := random(0, 11)
	return tacoQuotes[r]
}

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

// To allow for marshaling/unmarshaling an AttachmentField
type attachmentFieldJSON struct {
	Title     string `json:"title"`
	Value     string `json:"value"`
	Short     bool   `json:"short"`
	TitleLink string `json:"title_link"`
}

// A slashcommand response attachment
type AttachmentField struct {
	title     string `json:"title"`
	value     string `json:"value"`
	short     bool   `json:"short"`
	titleLink string `json:"title_link"`
}

func (a AttachmentField) Title() string     { return a.title }
func (a AttachmentField) Value() string     { return a.value }
func (a AttachmentField) Short() bool       { return a.short }
func (a AttachmentField) TitleLink() string { return a.titleLink }

// Marshal an AttachmentField using attachmentFieldJSON
func (a *AttachmentField) MarshalJSON() ([]byte, error) {
	marshalable := attachmentFieldJSON{
		Title:     a.title,
		Value:     a.value,
		Short:     a.short,
		TitleLink: a.titleLink}
	return json.Marshal(marshalable)
}

// Unmarshal an AttachmentField using attachmentFieldJSON
func (a *AttachmentField) UnmarshalJSON(data []byte) error {
	unmarshalable := attachmentFieldJSON{}
	err := json.Unmarshal(data, &unmarshalable)
	if err != nil {
		return err
	}
	a.title = unmarshalable.Title
	a.value = unmarshalable.Value
	a.short = unmarshalable.Short
	a.titleLink = unmarshalable.TitleLink

	return nil
}

// Verifies the slack command's token by matching it up to
// the environment variable TACOFANCY_VERIFICATION_TOKEN
// Returns true if the tokens match, otherwise returns false
func VerifyToken(token string) bool {
	return token == os.Getenv("TACOFANCY_VERIFICATION_TOKEN")
}

// Functional approach like this?
// or should I take a more C-like approach and pass in attachments and modify it
// or should this actually be a method of SlashCommand?
func BuildAttachments(taco tacofancy.Taco) []map[string]interface{} {
	attachments := make([]map[string]interface{}, 1)
	attachments[0] = make(map[string]interface{})
	fields := make([]AttachmentField, 5, 5)

	// FullTaco specific
	//attachments[0]["title"] = taco.Name()
	//attachments[0]["title_link"] = githubRawUrlToRepo(taco.URL)
	attachments[0]["text"] = taco.Description()
	if taco.BaseLayer().Name != "" {
		fields[0] = NewRecipeAttachmentField("Base Layer", taco.BaseLayer())
	}
	if taco.Seasoning().Name != "" {
		fields[1] = NewRecipeAttachmentField("Seasoning", taco.Seasoning())
	}
	if taco.Mixin().Name != "" {
		fields[2] = NewRecipeAttachmentField("Mixin", taco.Mixin())
	}
	if taco.Condiment().Name != "" {
		fields[3] = NewRecipeAttachmentField("Condiment", taco.Condiment())
	}
	if taco.Shell().Name != "" {
		fields[4] = NewRecipeAttachmentField("Shell", taco.Shell())
	}
	attachments[0]["fields"] = fields

	return attachments
}

func NewRecipeAttachmentField(title string, part tacofancy.TacoPart) AttachmentField {
	return AttachmentField{
		title: title + ": ",
		value: "<" + githubRawUrlToRepo(part.URL) + "|" + part.Name + ">",
		short: true}
}

// TODO: take the taco as an arg?  Then these become super easy to test since there is no api call
// Returns a SlashCommandResponse for the command `/taco recipe`
func NewTacoRecipeResponse(client tacofancy.Client) (SlashCommandResponse, error) {
	fullTaco, err := client.GetRandomFullTaco()
	if err != nil {
		return SlashCommandResponse{}, err
	}

	// TODO: or should I pass a reference, OO style?  I am not modifying it.
	attachments := BuildAttachments(fullTaco)
	// A bunch of duplicated stuff which it seems like using Taco interface would solve
	// but since Taco interface cannot access the struct properties, here we are... unless the two objects
	// just become one with some unused parts or I go with a whole bunch of duplicated getters/setters
	attachments[0]["title"] = fullTaco.Name()
	attachments[0]["title_link"] = githubRawUrlToRepo(fullTaco.URL())
	return NewSlashCommandResponse("in_channel", "", attachments), nil
}

// Returns a SlashCommandResponse for the command `/taco loco`
func NewTacoLocoResponse(client tacofancy.Client) (SlashCommandResponse, error) {
	randomTaco, err := client.GetRandomTacoParts()
	if err != nil {
		return SlashCommandResponse{}, err
	}

	attachments := BuildAttachments(randomTaco)
	attachments[0]["title"] = "A Delicious Random Taco"
	return NewSlashCommandResponse("in_channel", "", attachments), nil
}

func NewTacoGrandeResponse() (SlashCommandResponse, error) {
	attachments := make([]map[string]interface{}, 1)
	attachments[0] = make(map[string]interface{})
	attachments[0]["title"] = "Taco Grande"
	attachments[0]["title_link"] = "https://www.youtube.com/watch?v=mX18yNwqnMg"
	attachments[0]["text"] = getQuote()
	return NewSlashCommandResponse("in_channel", "", attachments), nil
}

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

// Creates and returns new SlashCommandResponse
func NewSlashCommandResponse(responseType, text string, attachments []map[string]interface{}) SlashCommandResponse {
	return SlashCommandResponse{responseType: responseType, text: text, attachments: attachments}
}

// Struct used for marshaling and unmarshaling a slashCommandResponseJSON
// want to call this a slashCommandResponseMarshaler, but the -er ending denotes
// an interface
type slashCommandResponseJSON struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
	// attachments could get more complex
	// https://api.slack.com/docs/message-attachments
	// using interface as the value data type here because it could be strings or a list
	Attachments []map[string]interface{} `json:"attachments"`
}

type SlashCommandResponse struct {
	responseType string `json:"response_type"`
	text         string `json:"text"`
	// attachments could get more complex
	// https://api.slack.com/docs/message-attachments
	// using interface as the value data type here because it could be strings or a list
	attachments []map[string]interface{} `json:"attachments"`
}

// Marshal an SlashCommandResponse using slashCommandResponseJSON
func (scr *SlashCommandResponse) MarshalJSON() ([]byte, error) {
	marshalable := slashCommandResponseJSON{
		ResponseType: scr.responseType,
		Text:         scr.text,
		Attachments:  scr.attachments}
	return json.Marshal(marshalable)
}

// Unmarshal an SlashCommandResponse using slashCommandResponseJSON
func (scr *SlashCommandResponse) UnmarshalJSON(data []byte) error {
	// Unmarshalable is confusing. It sounds like it cannot be marshalled, but
	// really it can be unmarshalled.  Unmarshaler would make sense, but again,
	// the -er ending is used to identify an interface.  Naming things is hard.
	unmarshalable := slashCommandResponseJSON{}
	err := json.Unmarshal(data, &unmarshalable)
	if err != nil {
		return err
	}
	scr.responseType = unmarshalable.ResponseType
	scr.text = unmarshalable.Text
	scr.attachments = unmarshalable.Attachments

	return nil
}

// Sends the provided SlashCommandResponse to the provided SlashCommand
// by Marshalling the SlashCommandResponse and making an HTTP POST to the SlashCommand.ResponseURL
func SendDelayedResponse(url string, scr SlashCommandResponse) error {
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
	commandResponseJson, err := json.Marshal(scr)
	if err != nil {
		return err
	}

	// bytes.NewBuffer() from https://stackoverflow.com/a/24455606
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(commandResponseJson))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
