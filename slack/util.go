package slack

import (
	"os"
	"strings"
)

// Verifies the slack command's token by matching it up to
// the environment variable TACOFANCY_VERIFICATION_TOKEN
// Returns true if the tokens match, otherwise returns false
func VerifyToken(token string) bool {
	return token == os.Getenv("TACOFANCY_VERIFICATION_TOKEN")
}

// converts https://raw.githubusercontent.com/sinker/tacofancy/master/condiments/baja_white_sauce.md
// to https://github.com/sinker/tacofancy/blob/master/condiments/baja_white_sauce.md
func githubRawUrlToRepo(url string) string {
	url = strings.Replace(url, "raw.github", "github", -1)
	url = strings.Replace(url, "master", "blob/master", -1)
	url = strings.Replace(url, "master//", "master/", -1)
	return url
}
