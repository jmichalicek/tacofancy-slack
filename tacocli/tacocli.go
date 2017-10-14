package main

import (
    "github.com/jmichalicek/tacofancy-slack/tacofancy"
    "fmt"
    "encoding/json"
    "flag"
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
    client := tacoFancy.NewClient("", nil)
    if *randomTaco {
        // showRandomTacoParts(*showJSON, *showDesc)
        taco, err := client.GetRandomTacoParts()
        if err != nil {
            fmt.Println(err)
            return
        }
        showTaco(&taco, *showJSON, *showDesc)
    } else {
        //showFullRandomTaco(*showJSON, *showDesc)
        taco, err := client.GetRandomFullTaco()
        if err != nil {
            fmt.Println(err)
            return
        }
        showTaco(&taco, *showJSON, *showDesc)
    }
}
