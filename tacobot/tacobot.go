package main

import (
    "net/http"
    "github.com/jmichalicek/tacofancy-slack/webhooks"
)


func main() {
	http.HandleFunc("/slack/slashcommand/", webhooks.SlashCommandHandler)
	http.ListenAndServe(":8080", nil)
}
