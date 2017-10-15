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

func getQuote() string {
	tacoQuotes := make([]string, 12, 12)
	tacoQuotes[0] = "Yo quiero chimichangas y chile colorado\nYo tengo el dinero para un steak picado\nLas flautas y tamales, siempre muy bueno\nY el chile relleno"
	tacoQuotes[1] = "You see, I just gotta have a tostada, carne asada\nThat's right, I want the whole enchilada\nMy only addiction has to do with a flour tortilla\nI need a quesadilla"
	tacoQuotes[2] = "I love to stuff my face with tacos al carbón\nWith my friends, or when I'm all alone\nYo tengo mucho hambre y ahora lo quiero\nUn burrito ranchero"
	tacoQuotes[3] = "So give me something spicy and hot, now\nBreak out the menu, what you got, now?\nOh, would you tell the waiter I'd like to have sour cream on the side\nYou better make sure the beans are refried"
	tacoQuotes[4] = "Well, there's not a taco big enough for a man like me\nThat's why I order two or three\nLet me give you a tip, just try a nacho chip\nIt's really good with bean dip"
	tacoQuotes[5] = "I eat uno, dos, tres, quatro burritos\nPretty soon I can't fit in my Speedos\nWell, I hope they feed us lots of chicken fajitas\nAnd a pitcher of margaritas"
	tacoQuotes[6] = "Well, the combination plates all come with beans and rice\nThe taquitos here are very nice\nNow I'm down on my knees, we need some extra tomatoes and cheese\nAnd could you make that separate checks, please?"
	tacoQuotes[7] = "Well, the food is coming, I can hardly wait\nNow watch your fingers, careful hot plate!\nWhat you think you're doing with my chile con queso?\nWell, if you want some, just say so"
	tacoQuotes[8] = "Oh boy, pico de gallo\nThey sure don't make it like this in Ohio\nNo gracias, yo quiero jalapeños, nada más\nYou can toss away the hot sauce"
	tacoQuotes[9] = "Donde estan los nachos? Holy frijole!\nYou better get me a bowl of guacamole\nY Usted, Eugene? Why's your face turning green?\nDon't you like pinto beans?"
	tacoQuotes[10] = "You want some more cinnamon crispas?\nIf you don't, hasta la vista\nJust take the rest home in a doggie bag if you wanna\nYou can finish it mañana"
	tacoQuotes[11] = "Well, it's been a pleasure, I can't eat no more\nSeñor, la cuenta, por favor\nIf you ain't ever tried real Mexican cooking, well, you oughta\nJust don't drink the water"

	random := func(min, max int) int {
		rand.Seed(time.Now().Unix())
		return rand.Intn(max-min) + min
	}

	r := random(0, 11)
	return tacoQuotes[r]
}

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

	// I do not really feel like this belongs here, but for the sake of easy testability
	// it will live here for now.  In a more complex app perhaps there would be a
	// TacoFancySlashCommand where this would feel more appropriate
	// or perhaps BuildResponse() should not be a SlashCommand method and I should pass
	// the SlashCommand and client into NewTacoRecipeResponse() etc then a factory method
	// which is basiclaly what BuildResponse() is.  The problem here is that I feel like
	// no matter what, at some point, I am configuring a tacofancy.Client and passing it inot
	// a function, not because I know the actual code needs it, but because the code "might" need it
	// and test cases need some way to pass it in for those.
	TacofancyClient	tacofancy.Client
}

// A slashcommand response attachment
type AttachmentField struct {
	Title     string `json:"title"`
	Value     string `json:"value"`
	Short     bool   `json:"short"`
	TitleLink string `json:"title_link"`
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
	baseLayer := taco.BaseLayer()
	if baseLayer.Name != "" {
		fields[0] = AttachmentField{Title: "Base Layer: ", Value: "<" + githubRawUrlToRepo(baseLayer.URL) + "|" + baseLayer.Name + ">", Short: true}
	}
	if taco.Seasoning().Name != "" {
		fields[1] = AttachmentField{Title: "Seasoning: ", Value: "<" + githubRawUrlToRepo(taco.Seasoning().URL) + "|" + taco.Seasoning().Name + ">", Short: true}
	}
	if taco.Mixin().Name != "" {
		fields[2] = AttachmentField{Title: "Mixin: ", Value: "<" + githubRawUrlToRepo(taco.Mixin().URL) + "|" + taco.Mixin().Name + ">", Short: true}
	}
	if taco.Condiment().Name != "" {
		fields[3] = AttachmentField{Title: "Condiment: ", Value: "<" + githubRawUrlToRepo(taco.Condiment().URL) + "|" + taco.Condiment().Name + ">", Short: true}
	}
	if taco.Shell().Name != "" {
		fields[4] = AttachmentField{Title: "Shell: ", Value: "<" + githubRawUrlToRepo(taco.Shell().URL) + "|" + taco.Shell().Name + ">", Short: true}
	}
	attachments[0]["fields"] = fields

	return attachments
}

func NewRecipeAttachmentField(title string, part tacofancy.TacoPart) AttachmentField {
	return AttachmentField{
		Title: title+": ",
		Value: "<" + githubRawUrlToRepo(part.URL) + "|" + part.Name + ">",
		Short: true}
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
	return SlashCommandResponse{ResponseType: "in_channel", Text: "", Attachments: attachments}, nil
}

// Returns a SlashCommandResponse for the command `/taco loco`
func NewTacoLocoResponse(client tacofancy.Client) (SlashCommandResponse, error) {
	randomTaco, err := client.GetRandomTacoParts()
	if err != nil {
		return SlashCommandResponse{}, err
	}

	attachments := BuildAttachments(randomTaco)
	attachments[0]["title"] = "A Delicious Random Taco"
	return SlashCommandResponse{ResponseType: "in_channel", Text: "", Attachments: attachments}, nil
}

func NewTacoGrandeResponse() (SlashCommandResponse, error) {
	attachments := make([]map[string]interface{}, 1)
	attachments[0] = make(map[string]interface{})
	// fields := make([]AttachmentField, 5, 5)

	attachments[0]["title"] = "Taco Grande"
	attachments[0]["title_link"] = "https://www.youtube.com/watch?v=mX18yNwqnMg"
	attachments[0]["text"] = getQuote()
	return SlashCommandResponse{ResponseType: "in_channel", Text: "", Attachments: attachments}, nil
}

// Builds a SlashCommandResponse to use for responding to a Slack SlashCommand
// TODO:  make a responder interface and pass that in?  then t.Respond() could
// be called?  Not sure since that moves the slack specific logic out of here
// but could almost certainly be far more useful for a more general slack framework
func (sc *SlashCommand) BuildResponse() (SlashCommandResponse, error) {
	// respond to slash command
	// TODO: Have a response factory registry?
	parts := strings.Split(sc.Text, " ")
	commandType := "recipe"
	if len(parts) >= 1 && parts[0] != "" {
		commandType = strings.ToLower(parts[0])
	}

	if commandType == "recipe" {
		return NewTacoRecipeResponse(sc.TacofancyClient)
	} else if commandType == "loco" {
		return NewTacoLocoResponse(sc.TacofancyClient)
	} else if commandType == "grande" {
		return NewTacoGrandeResponse()
	}
	return SlashCommandResponse{ResponseType: "in_channel", Text: "I didn't understand that command"}, nil
}

// Builds a SlashCommandResponse and calls SendDelayedResponse() and returns the SendDelayedResponse()'s
// return value in the channel.
func (sc *SlashCommand) RespondAsync(ch chan error) {
	defer close(ch)
	scr, _ := sc.BuildResponse()
	ch <- SendDelayedResponse(sc.ResponseURL, scr)
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
