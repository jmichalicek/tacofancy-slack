package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jmichalicek/tacofancy-slack/tacofancy"
)

func showTaco(taco tacofancy.Taco, showJSON bool, showDesc bool) {
	s, err := json.MarshalIndent(taco, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	if showJSON {
		fmt.Println(string(s))
	}

	if showDesc {
		fmt.Println(taco.Description())
	}
}

func main() {
	showDesc := flag.Bool("description", false, "Show the description output")
	showJSON := flag.Bool("json", false, "Show the JSON output")
	randomTaco := flag.Bool(
		"random", false, "Show a random selection of base layer, seasoning, mixin, condiment, and shell")
	flag.Parse()

	// default url and http.Client
	client := tacofancy.NewClient("", nil)
	if *randomTaco {
		taco, err := client.GetRandomTacoParts()
		if err != nil {
			fmt.Println(err)
			return
		}
		showTaco(&taco, *showJSON, *showDesc)
	} else {
		taco, err := client.GetRandomFullTaco()
		if err != nil {
			fmt.Println(err)
			return
		}
		// ideally no pointer here, but the unmarshall seems to need it.
		// other option is type switch stuff in showTaco
		showTaco(&taco, *showJSON, *showDesc)
	}
}
