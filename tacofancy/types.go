package tacofancy

import "encoding/json"

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

	// Should these take pointers?
	SetBaseLayer(TacoPart)
	SetMixin(TacoPart)
	SetCondiment(TacoPart)
	SetShell(TacoPart)
	SetSeasoning(TacoPart)
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
	*baseTacoJSON
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

// Should these take pointers?  I assume that the struct would need to use
// pointers for that to matter?
func (t *BaseTaco) SetBaseLayer(bl TacoPart) { t.baseLayer = bl }
func (t *BaseTaco) SetMixin(m TacoPart)      { t.mixin = m }
func (t *BaseTaco) SetCondiment(c TacoPart)  { t.condiment = c }
func (t *BaseTaco) SetSeasoning(s TacoPart)  { t.seasoning = s }
func (t *BaseTaco) SetShell(s TacoPart)      { t.shell = s }

func NewBaseTaco(baseLayer, mixin, condiment, seasoning, shell TacoPart) BaseTaco {
	return BaseTaco{baseLayer: baseLayer, mixin: mixin, condiment: condiment, seasoning: seasoning, shell: shell}
}

// A taco which is not from a recipe, but from random selection of parts
type RandomTaco struct {
	BaseTaco
}

// Returns a description of the taco made up from the names of the parts
func (t *RandomTaco) Description() string {
	// shell names are inconsistent, but roll with this for now
	desc := t.baseLayer.Name + " seasoned with " + t.seasoning.Name + " with " + t.mixin.Name + " and " +
		t.condiment.Name + " in " + t.shell.Name + "."
	return desc
}

// Implements Marshaler and Unmarshaler on RandomTaco
// to deal with the private fields
// Unmarshals a RandomTacoJSON and then uses it to call the setters
// on RandomTaco
func (t *RandomTaco) UnmarshalJSON(data []byte) error {
	rtj := baseTacoJSON{}
	err := json.Unmarshal(data, &rtj)
	if err != nil {
		return err
	}
	t.SetBaseLayer(rtj.BaseLayer)
	t.SetMixin(rtj.Mixin)
	t.SetSeasoning(rtj.Seasoning)
	t.SetCondiment(rtj.Condiment)
	t.SetShell(rtj.Shell)

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

// Implementation of a Taco from a full recipe
// Fields duplicated from RandomTaco because de-duplicating those plus
// allowing for json marshal/unmarshal plus using idiomatically named
// getters is a pain.  End up implementing getters and setters for each struct
// MarshalJSON and UnmarshalJSON for each struct, etc.  It would be less work
// to just have Taco and have some extra fields.
// Experimented with BaseTaco embedded struct which had the getters and setters
// May try it again.
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

// Ignoring the RecipeURL, etc. for now
func (t *FullTaco) Name() string       { return t.name }
func (t *FullTaco) URL() string        { return t.url }
func (t *FullTaco) Recipe() string     { return t.recipe }
func (t *FullTaco) Slug() string       { return t.slug }
func (t *FullTaco) SetName(n string)   { t.name = n }
func (t *FullTaco) SetURL(u string)    { t.url = u }
func (t *FullTaco) SetRecipe(r string) { t.recipe = r }
func (t *FullTaco) SetSlug(s string)   { t.slug = s }

// Returns a description of the taco made up from the names of the parts
func (t *FullTaco) Description() string {
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
func (t *FullTaco) UnmarshalJSON(data []byte) error {
	tj := fullTacoJSON{}
	err := json.Unmarshal(data, &tj)
	if err != nil {
		return err
	}
	t.SetBaseLayer(tj.BaseLayer)
	t.SetMixin(tj.Mixin)
	t.SetSeasoning(tj.Seasoning)
	t.SetCondiment(tj.Condiment)
	t.SetShell(tj.Shell)

	t.SetName(tj.Name)
	t.SetURL(tj.URL)
	t.SetRecipe(tj.Recipe)
	t.SetSlug(tj.Slug)

	return nil
}

// Creates a RandomTacoJSON, which is required by json.Marshal to access
// the public fields, and then returns the marshal data
func (t *FullTaco) MarshalJSON() ([]byte, error) {
	tj := fullTacoJSON{
		baseTacoJSON: &baseTacoJSON{BaseLayer: t.baseLayer,
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
