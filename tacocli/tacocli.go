package main

import (
    "github.com/jmichalicek/tacofancy-slack/tacofancy"
    "fmt"
    "encoding/json"
    "flag"
)

func showRandomTacoParts(showJSON, showDesc bool) {
    taco, err := tacofancy.GetRandomTacoParts()

    // move this back to main() by using a Taco interface
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

func showFullRandomTaco(showJSON, showDesc bool) {
    taco, err := tacofancy.GetRandomFullTaco()

    // move this back to main() by using a Taco interface
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

    // I could use a Taco interface to reduce duplication here, I think
    if *randomTaco {
        showRandomTacoParts(*showJSON, *showDesc)
    } else {
        showFullRandomTaco(*showJSON, *showDesc)
    }

    // s, err := json.MarshalIndent(taco, "", "  ")
    // if err != nil {
    //     fmt.Println(err)
    // }
    //
    // if *showJSON {
    //     fmt.Println(string(s))
    // }
    //
    // if *showDesc {
    //     fmt.Println(taco.Description())
    // }
}
