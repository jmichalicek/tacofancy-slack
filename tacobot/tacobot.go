package main

import (
    "net/http"
    "github.com/jmichalicek/tacofancy-slack/webhooks"
)


// Runs the web interface for the tacobot
func main() {
	http.HandleFunc("/slack/slashcommand/", webhooks.SlashCommandHandler)
	http.ListenAndServe(":8080", nil)
}
