package webhooks

import (
    "net/http"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
    // TODO: Parse the incoming json and do stuff
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(http. StatusNoContent)
}
