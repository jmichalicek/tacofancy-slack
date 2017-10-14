package slack

import (
	"encoding/json"
	"github.com/jmichalicek/tacofancy-slack/tacofancy"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestSendDelayedResponse(t *testing.T) {
	// TODO
}

var randomApiResponse string = `{
  "base_layer": {
    "url": "https://raw.github.com/sinker/tacofancy/master/base_layers/slow_cooked_salsa_verde_chicken.md",
    "name": "Slow-Cooked Salsa Verde Chicken",
    "slug": "slow_cooked_salsa_verde_chicken",
    "recipe": "Slow-Cooked Salsa Verde Chicken\n===============================\n\nLow-effort, delicious shredded chicken based off [this recipe](http://www.skinnytaste.com/2013/04/easiest-crock-pot-salsa-verde-chicken.html). To minimize prep time, just toss in a jar of storebought salsa verde.\n\n* 2 lbs skinless chicken breasts\n* 2 cups salsa verde\n* 1 tsp minced garlic or 1/4 tsp garlic powder\n* Pinch of Oregano _(Note: I ran out of Oregano, but the recipe still tastes great without it.)_\n* Pinch of Cumin\n* Salt, to taste\n\nAdd chicken to slow cooker and season with garlic, oregano, cumin and salt. Pour salsa verde over everything, cover and cook for two hours on high.\n\nOnce that's ready, shred chicken. Give it another stir to cover everything in sauce, then serve."
  },
  "seasoning": {
    "url": "https://raw.github.com/sinker/tacofancy/master/seasonings/taco_and_fajita.md",
    "name": "Kathy's Taco and Fajita Seasoning",
    "slug": "kathys_taco_and_fajita_seasoning",
    "recipe": "Kathy's Taco and Fajita Seasoning\n=================================\n\nMy cousin Kathy sent me this recipe for my birthday:\n\nAmazing taco/fajita seasoning – Use with any protein.  It's a great recipe to quadruple so you have it ready for immediate taco needs. When making a big batch, use less salt.\n\nCombine:\n\n2 tsp cumin\n1 tsp EACH: paprika, onion powder, dried oregano, and kosher salt\n½ tsp EACH: ground coriander, garlic powder and black pepper\n¼ tsp EACH: red pepper flakes, ground ginger and cayenne pepper\n\ntags: vegetarian, vegan\n"
  },
  "mixin": {
    "url": "https://raw.github.com/sinker/tacofancy/master/mixins/veg_for_fish_tacos.md",
    "name": "Veggies for Fish Tacos",
    "slug": "veggies_for_fish_tacos",
    "recipe": "Veggies for Fish Tacos\n======================\n\nFish tacos are a special breed, requiring different vegetable options.\n\n__Assemble your veg from the following options:__\n\n* Cabbage, purple, shredded\n* Cabbage, other shades, shredded\n* Radishes, sliced into thin slices\n* Red peppers, diced\n* Cherry tomatoes, sliced (if you're a heathen)\n* Cilantro, if it doesn't taste like soap to you\n\nAnd one requirement:\n* Limes, sliced for juicing over tacos.\n\nPlace out your selections and assemble into your taco. Then squeeze a lime over the top.\n\ntags: vegetarian, vegan\n"
  },
	"condiment": {
    "url": "https://raw.github.com/sinker/tacofancy/master/condiments/simple_salsa_verde.md",
    "name": "Simple Salsa Verde",
    "slug": "simple_salsa_verde",
    "recipe": "Simple Salsa Verde\n==================\n\nI got this base recipe from a vegan friend. If you can't find one of these peppers, swap in another one!\n\n* 6 Average-sized tomatillos\n* 1 Poblano pepper\n* 1 Serrano pepper\n* 1 Jalapeno pepper\n* 1 Sweet red pepper\n* Juice of 1 or 2 fresh-squeezed limes (to taste)\n* Pinch or two kosher salt (to taste)\n\nYou're in charge of the heat here. For a milder salsa, remove all the ribs and seeds inside the peppers. For medium, leave in a few ribs, and for hot, go nuts. Rough chop the peppers and tomatillos, then throw into a blender or food processor with salt and lime juice. Pulse to desired consistency.\n\nAs with most salsas, this will taste better if you let it sit in the fridge for a few hours before eating. It's great on chips or drizzled over steak or pork tacos.\n\ntags: vegetarian, vegan\n"
  },
  "shell": {
    "url": "https://raw.github.com/sinker/tacofancy/master/shells/Fresh_corn_tortillas.md",
    "name": "Fresh Corn Tortillas",
    "slug": "fresh_corn_tortillas",
    "recipe": "Fresh Corn Tortillas\n===================\n\nThis is the only way to go. So worth it. Makes roughly 15 tortillas.\n\n* 1 3/4 cups masa harina\n* 1 1/8 cups water\n\n1. In a medium bowl, mix together masa harina and hot water until thoroughly combined. Turn dough onto a clean surface and knead until pliable and smooth. If dough is too sticky, add more masa harina; if it begins to dry out, sprinkle with water. Cover dough tightly with plastic wrap and allow to stand for 30 minutes.\n2. Preheat a cast iron skillet or griddle to medium-high.\n3. Divide dough into 15 equal-size balls. Using a tortilla press (or a rolling pin), press each ball of dough flat between two sheets of wax paper (plastic wrap or a freezer bag cut into halves will also work).\n4. Place tortilla in preheated pan and allow to cook for approximately 30 seconds, or until browned and slightly puffy. Turn tortilla over to brown on second side for approximately 30 seconds more, then transfer to a plate. Repeat process with each ball of dough. Keep tortillas covered with a towel to stay warm and moist (or a low temp oven) until ready to serve.\n\ntags: vegetarian, vegan\n"
  }
}`

var randomApiData tacofancy.RandomTaco = tacofancy.RandomTaco{}
var recipeApiData tacofancy.FullTaco = tacofancy.FullTaco{}

func TestSlashCommandBuildResponseForTacoLoco(t *testing.T) {

	var apiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(randomApiResponse))
	}))
	client := tacofancy.NewClient(apiStub.URL, nil)
	slashCommand := SlashCommand{
		Command: "/taco", Text: "loco",
		Token: "testtoken", ResponseURL: "https://example.org/",
		ChannelId: "fake", UserId: "fake", TacofancyClient: client}

	response, err := slashCommand.BuildResponse()

	if err != nil {
		t.Errorf("Got error: ", err)
	}

	json.Unmarshal([]byte(randomApiResponse), &randomApiData)
	attachments := BuildAttachments(&randomApiData) // from slack.go
	attachments[0]["title"] = "A Delicious Random Taco"
	expectedResponse := SlashCommandResponse{
		ResponseType: "in_channel",
		Text:         "",
		Attachments:  attachments}

	if expectedResponse.ResponseType != response.ResponseType {
		t.Errorf("oh no")
	}

	if expectedResponse.Text != response.Text {
		t.Errorf("oh no")
	}

	if !reflect.DeepEqual(response.Attachments[0], expectedResponse.Attachments[0]) {
		t.Errorf("Expected SlashCommandResponse of:  %v\n\n but got: %v", expectedResponse, response)
	}
}

func TestSlashCommandRespondAsync(t *testing.T) {
	// TODO
}

func TestVerifyToken(t *testing.T) {
	// TODO
}

func TestBuildAttachments(t *testing.T) {
	
}
