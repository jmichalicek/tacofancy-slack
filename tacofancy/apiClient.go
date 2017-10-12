package tacofancy

import (
	"net/http"
	//"github.com/jmichalicek/tacofancy-slack/tacofancy"
	"encoding/json"
	"io/ioutil"
	"time"
)

const baseUrl string = "https://taco-randomizer.herokuapp.com/"
const randomPath string = "random/"

var BaseUrl = baseUrl

// default client has no timeout, so we make our own with a timeout
var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

func GetRandomTacoParts(client *http.Client) (RandomTaco, error) {
	// random base layer, mixin, condiment, seasoning, and shell
	if client == nil {
		client = httpClient
	}
	var tacoUrl = baseUrl + randomPath
	var taco = RandomTaco{}

	r, err := client.Get(tacoUrl)
	if err != nil {
		return taco, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal(body, &taco)
	return taco, err
}

func GetRandomFullTaco(client *http.Client) (FullTaco, error) {
	if client == nil {
		client = httpClient
	}
	// there's probably a built in url path  manipulation thing in Go, but I am being lazy for now
	var tacoUrl = baseUrl + randomPath + "?full-taco=true"
	var taco = FullTaco{}

	r, err := client.Get(tacoUrl)
	if err != nil {
		return taco, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal(body, &taco)
	return taco, err
}
