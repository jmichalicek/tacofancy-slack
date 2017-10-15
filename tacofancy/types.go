package tacofancy

import (
	"encoding/json"
)

// Some serious over-engineering is going on here for the sake of learning Go.
// There is a Taco interface, BaseTaco which implements it,
// RandomTaco and FullTaco which embed BaseTaco so that getters/setters/attributes
// do not have to be duplicated.
// And then baseTacoJSON and fullTacoJSON which have exported attributes, but
// those structs themselves are not exported which allows FullTaco and
// RandomTaco to Marshal and Unmarshal JSON while having all unexported
// member attributes which are ignored by json.Marshal() and json.Unmarhsal()
// This mess with the *JSON structs duplicating fields is perhaps my least
// favorite thing so far with Go's take on OOP

// An interface implementing a tco
type Taco interface {
	Description() string
	BaseLayer() TacoPart
	Mixin() TacoPart
	Condiment() TacoPart
	Seasoning() TacoPart
	Shell() TacoPart
}

// see https://taco-randomizer.herokuapp.com/random/ for a random example
// Seems to be consistent in one base layer one shell, one seasoning,
// one mixin, one condiment

// Implements one igredient or part from a tacofancy taco
// A TacoPart is made up of a Name, URL, Slug, and Recipe
type TacoPart struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Slug string `json:"slug"`
	// Recipe instructions on how to make this TacoPart
	Recipe string `json:"recipe"`

	// could make sense to put the type of part on here, but the api does not return that data
}

func NewTacoPart(name, url, slug, recipe string) TacoPart {
	return TacoPart{Name: name, URL: url, Slug: slug, Recipe: recipe}
}

// FOr marshaling and unmarshaling a BaseTaco
type baseTacoJSON struct {
	// The actual baselayer
	BaseLayer TacoPart `json:"base_layer"`
	// The actual baselayer
	Mixin TacoPart `json:"mixin"`
	// The actual baselayer
	Condiment TacoPart `json:"condiment"`
	// The actual baselayer
	Seasoning TacoPart `json:"seasoning"`
	// The actual baselayer
	Shell TacoPart `json:"shell"`
}

type fullTacoJSON struct {
	baseTacoJSON
	Name string `json:"name"`
	// The URL to the recipe
	URL string `json:"url"`
	// Recipe instructions on how to make this Taco
	Recipe string `json:"recipe"`
	// Slug part of the URL
	Slug string `json:"slug"`

	// full taco api doubles up url stuff for some reaosn
	// URL to the baselayer for the taco
	BaseLayerURL string `json:"base_layer_url"`
	// URL to the mixin of the taco
	MixinURL string `json:"mixin_url"`
	// URL to the condiment of the taco
	CondimentURL string `json:"condiment_url"`
	// URL to the seasoning of the taco
	SeasoningURL string `json:"seasoning_url"`
	// URL to the shell of the taco
	ShellURL string `json:"shell_url"`
}

func NewBaseTaco(baseLayer, mixin, condiment, seasoning, shell TacoPart) BaseTaco {
	return BaseTaco{baseLayer: baseLayer, mixin: mixin, condiment: condiment, seasoning: seasoning, shell: shell}
}

// Implementation of a Taco made of random parts
// This has a BaseLayer, Mixin, Condiment, Seasoning, and Shell
// TacoPart
type BaseTaco struct {
	// The actual baselayer
	baseLayer TacoPart `json:"base_layer"`
	// The actual baselayer
	mixin TacoPart `json:"mixin"`
	// The actual baselayer
	condiment TacoPart `json:"condiment"`
	// The actual baselayer
	seasoning TacoPart `json:"seasoning"`
	// The actual baselayer
	shell TacoPart `json:"shell"`
}

func (t BaseTaco) Description() string {
	// shell names are inconsistent, but roll with this for now
	desc := t.baseLayer.Name + " seasoned with " + t.seasoning.Name + " with " + t.mixin.Name + " and " +
		t.condiment.Name + " in " + t.shell.Name + "."
	return desc
}

// Should these return pointers?
func (t BaseTaco) BaseLayer() TacoPart { return t.baseLayer }
func (t BaseTaco) Mixin() TacoPart     { return t.mixin }
func (t BaseTaco) Condiment() TacoPart { return t.condiment }
func (t BaseTaco) Seasoning() TacoPart { return t.seasoning }
func (t BaseTaco) Shell() TacoPart     { return t.shell }

// Create a new randomly generated taco
// These tacos are made of random parts as returned from the TacoFancy API.
func NewRandomTaco(baseLayer, mixin, condiment, seasoning, shell TacoPart) RandomTaco {
	return RandomTaco{
		BaseTaco: BaseTaco{baseLayer: baseLayer, mixin: mixin, condiment: condiment, seasoning: seasoning, shell: shell}}
}

// A taco which is not from a recipe, but from random selection of parts
type RandomTaco struct {
	BaseTaco
}

// Returns a description of the taco made up from the names of the parts
func (t RandomTaco) Description() string {
	// shell names are inconsistent, but roll with this for now
	desc := t.baseLayer.Name + " seasoned with " + t.seasoning.Name + " with " + t.mixin.Name + " and " +
		t.condiment.Name + " in " + t.shell.Name + "."
	return desc
}

// Implements Marshaler and Unmarshaler on RandomTaco
// to deal with the private fields
// Unmarshals a RandomTacoJSON and then uses it to call the setters
// on RandomTaco
// This is bending the rules a bit on immutability, maybe there's a better answer
// but to do it right, we'd need to return the RandomTaco.  Perhaps a "proper" way would be
// to wrap this in some manner to make a new taco, unmarshal into it, and then return that
func (t *RandomTaco) UnmarshalJSON(data []byte) error {
	rtj := baseTacoJSON{}
	err := json.Unmarshal(data, &rtj)
	if err != nil {
		return err
	}
	t.baseLayer = rtj.BaseLayer
	t.mixin = rtj.Mixin
	t.seasoning = rtj.Seasoning
	t.condiment = rtj.Condiment
	t.shell = rtj.Shell

	return nil
}

// Creates a RandomTacoJSON, which is required by json.Marshal to access
// the public fields, and then returns the marshal data
func (t *RandomTaco) MarshalJSON() ([]byte, error) {
	rtj := baseTacoJSON{
		BaseLayer: t.baseLayer,
		Mixin:     t.mixin,
		Seasoning: t.seasoning,
		Condiment: t.condiment,
		Shell:     t.shell}
	return json.Marshal(rtj)
}

// Create a new randomly generated taco
// These tacos are made of random parts as returned from the TacoFancy API.
func NewFullTaco(name, url, recipie, slug string, baseLayer, mixin, condiment, seasoning, shell TacoPart) FullTaco {
	return FullTaco{
		BaseTaco: BaseTaco{baseLayer: baseLayer, mixin: mixin, condiment: condiment, seasoning: seasoning, shell: shell}}
}

// Implementation of a Taco from a full recipe
// These tacos come from the TacoFancy API as a complete recipe with a name, url, slug, etc.
type FullTaco struct {
	// Can also do this as *BaseTaco but I am not clear on the implications
	// of using a pointer and it seems to be breaking setters.
	BaseTaco
	// The name of the taco
	name string `json:"name"`
	// The URL to the recipe
	url string `json:"url"`
	// Recipe instructions on how to make this Taco
	recipe string `json:"recipe"`
	// Slug part of the URL
	slug string `json:"slug"`

	// full taco api doubles up url stuff for some reaosn
	// URL to the baselayer for the taco
	baseLayerURL string `json:"base_layer_url"`
	// URL to the mixin of the taco
	mixinURL string `json:"mixin_url"`
	// URL to the condiment of the taco
	condimentURL string `json:"condiment_url"`
	// URL to the seasoning of the taco
	seasoningURL string `json:"seasoning_url"`
	// URL to the shell of the taco
	shellURL string `json:"shell_url"`
}

// Ignoring the RecipeURL, etc. for now
func (t FullTaco) Name() string         { return t.name }
func (t FullTaco) URL() string          { return t.url }
func (t FullTaco) Recipe() string       { return t.recipe }
func (t FullTaco) Slug() string         { return t.slug }
func (t FullTaco) BaseLayerURL() string { return t.baseLayerURL }
func (t FullTaco) MixinURL() string     { return t.mixinURL }
func (t FullTaco) CondimentURL() string { return t.condimentURL }
func (t FullTaco) SeasoningURL() string { return t.seasoningURL }
func (t FullTaco) ShellURL() string     { return t.shellURL }

// Returns a description of the taco made up from the names of the parts
func (t FullTaco) Description() string {
	desc := t.baseLayer.Name
	if t.seasoning.Name != "" {
		desc = desc + " seasoned with " + t.seasoning.Name
	}

	if t.mixin.Name != "" {
		desc = desc + " with " + t.mixin.Name
	}

	if t.condiment.Name != "" {
		joinWord := " with "
		if t.mixin.Name != "" {
			joinWord = " and "
		}
		desc = desc + joinWord + t.condiment.Name
	}

	// full tacos do not seem to have a shell usually and shell names are
	// inconsistent
	if t.shell.Name != "" {
		desc = desc + " in " + t.shell.Name + "."
	}

	return desc
}

// Implements Marshaler and Unmarshaler on RandomTaco
// to deal with the private fields
// Unmarshals a RandomTacoJSON and then uses it to call the setters
// on RandomTaco
// mucks up the whole immutability thing for now
func (t *FullTaco) UnmarshalJSON(data []byte) error {
	tj := fullTacoJSON{}
	err := json.Unmarshal(data, &tj)
	if err != nil {
		return err
	}
	t.baseLayer = tj.BaseLayer
	t.mixin = tj.Mixin
	t.seasoning = tj.Seasoning
	t.condiment = tj.Condiment
	t.shell = tj.Shell
	t.name = tj.Name
	t.url = tj.URL
	t.recipe = tj.Recipe
	t.slug = tj.Slug
	t.baseLayerURL = tj.BaseLayerURL
	t.mixinURL = tj.MixinURL
	t.condimentURL = tj.CondimentURL
	t.seasoningURL = tj.SeasoningURL
	t.shellURL = tj.ShellURL

	return nil
}

// Creates a RandomTacoJSON, which is required by json.Marshal to access
// the public fields, and then returns the marshal data
func (t *FullTaco) MarshalJSON() ([]byte, error) {
	tj := fullTacoJSON{
		baseTacoJSON: baseTacoJSON{BaseLayer: t.baseLayer,
			Mixin:     t.mixin,
			Seasoning: t.seasoning,
			Condiment: t.condiment,
			Shell:     t.shell},
		Name:   t.name,
		URL:    t.url,
		Recipe: t.recipe,
		Slug:   t.slug}
	return json.Marshal(tj)
}
