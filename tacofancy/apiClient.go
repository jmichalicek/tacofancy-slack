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

// default client has no timeout, so we make our own with a timeout
var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

func GetRandomTacoParts() (RandomTaco, error) {
	// random base layer, mixin, condiment, seasoning, and shell
	var tacoUrl = baseUrl + randomPath
	var taco = RandomTaco{}

	r, err := httpClient.Get(tacoUrl)
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

func GetRandomFullTaco() (FullTaco, error) {
	// there's probably a built in url path  manipulation thing in Go, but I am being lazy for now
	var tacoUrl = baseUrl + randomPath + "?full-taco=true"
	var taco = FullTaco{}

	r, err := httpClient.Get(tacoUrl)
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
