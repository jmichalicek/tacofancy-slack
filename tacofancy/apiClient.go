package tacofancy

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"time"
)

const DefaultBaseURL string = "https://taco-randomizer.herokuapp.com"
const randomPath string = "/random/"
// using reference to http.Client here like https://github.com/digitalocean/godo/blob/master/godo.go#L151
// and the golang github client written by google
var DefaultHTTPClient = &http.Client{Timeout: time.Second * 10}

// Get a new tacofancy api client
// should this return a reference like the official golang github client
// written by google does as well as the digital ocean one?
// https://github.com/digitalocean/godo/blob/master/godo.go#L151
func NewClient(baseURL string, httpClient *http.Client) Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	if httpClient == nil {
		httpClient = DefaultHTTPClient
	}

	return Client{BaseURL: baseURL, httpClient: httpClient}
}

type Client struct {
  httpClient *http.Client
  BaseURL string
}

func (client Client) GetRandomTacoParts() (RandomTaco, error) {
	// random base layer, mixin, condiment, seasoning, and shell

	var tacoUrl = client.BaseURL + randomPath
	var taco = RandomTaco{}

	r, err := client.httpClient.Get(tacoUrl)
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

func (client Client) GetRandomFullTaco() (FullTaco, error) {
	// there's probably a built in url path  manipulation thing in Go, but I am being lazy for now
	var tacoUrl = client.BaseURL + randomPath + "?full-taco=true"
	var taco = FullTaco{}

	r, err := client.httpClient.Get(tacoUrl)
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
