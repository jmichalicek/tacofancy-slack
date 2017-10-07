package tacofancy


// An interface implementing a tco
type Taco interface {
    Description() string
    BaseLayer() *TacoPart
    SetBaseLayer(*TacoPart)
    Mixin() *TacoPart
    SetMixin(*TacoPart)
    Condiment() *TacoPart
    SetCondiment(*TacoPart)
    Seasoning() *TacoPart
    SetSeasoning(*TacoPart)
    Shell() *TacoPart
    SetShell(*TacoPart)
}

// see https://taco-randomizer.herokuapp.com/random/ for a random example
// Seems to be consistent in one base layer one shell, one seasoning,
// one mixin, one condiment

// Implements one igredient or part from a tacofancy taco
// A TacoPart is made up of a Name, URL, Slug, and Recipe
type TacoPart struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Slug   string `json:"slug"`
    // Recipe instructions on how to make this TacoPart
	Recipe string `json:"recipe"`
}

// Implements the common parts and methods of a Taco
type CommonTaco struct {
    // The actual baselayer
	baseLayer TacoPart `json:"base_layer"`
    // The actual baselayer
	mixin     TacoPart `json:"mixin"`
    // The actual baselayer
	condiment TacoPart `json:"condiment"`
    // The actual baselayer
	seasoning TacoPart `json:"seasoning"`
    // The actual baselayer
	shell     TacoPart `json:"shell"`
}

func (t CommonTaco) BaseLayer() *TacoPart {
    return &t.baseLayer
}

func (t CommonTaco) Mixin() *TacoPart {
    return &t.mixin
}

func (t CommonTaco) Condiment() *TacoPart {
    return &t.condiment
}

func (t CommonTaco) Seasoning() *TacoPart {
    return &t.seasoning
}

func (t CommonTaco) Shell() *TacoPart {
    return &t.shell
}

func (t *CommonTaco) SetBaseLayer(baseLayer *TacoPart)  {
    t.baseLayer = *baseLayer
}

func (t *CommonTaco) SetMixin(mixin *TacoPart)  {
    t.mixin = *mixin
}

func (t *CommonTaco) SetCondiment(condiment *TacoPart)  {
    t.condiment = *condiment
}

func (t *CommonTaco) SetSeasoning(seasoning *TacoPart)  {
    t.seasoning = *seasoning
}

func (t *CommonTaco) SetShell(shell *TacoPart)  {
    t.shell = shell
}

// Implementation of a Taco made of random parts
// This has a BaseLayer, Mixin, Condiment, Seasoning, and Shell
// TacoPart
type RandomTaco struct {
    CommonTaco
}

// Returns a description of the taco made up from the names of the parts
func (t *RandomTaco) Description() string {

	// shell names are inconsistent, but roll with this for now
	_, desc := t.baseLayer.Name+" seasoned with "+t.seasoning.Name+" with "+t.mixin.Name+" and "+
		t.condiment.Name+" in "+t.shell.Name+"."
	return desc
}

func (t *RandomTaco) UnmarshalJSON() data []byte, error {
    data []byte
	doc := RandomTacoJSON{
        BaseLayer: t.baseLayer, Mixin: t.mixin, Condiment: t.condiment,
        Seasoning: t.seasoning, Shell: t.shell}
	err := json.Unmarshal(data, doc)
	if err != nil {
		return err
	}

	return data
}

// Provides a struct for marshaling and unmarshaling a RandomTaco
type RandomTacoJSON struct {
    BaseLayer    TacoPart `json:"base_layer"`
    // The mixin
    Mixin        TacoPart `json:"mixin"`
    Condiment    TacoPart `json:"condiment"`
    // The seasoning
    Seasoning    TacoPart `json:"seasoning"`
    // The shell
    Shell        TacoPart `json:"shell"`
}

// Marshal's to a RandomTaco
func (t *RandomTacoJSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Implementation of a Taco from a full recipe
type FullTaco struct {
    // The name of the taco
	Name   string `json:"name"`
    // The URL to the recipe
	URL    string `json:"url"`
    // Recipe instructions on how to make this Taco
	Recipe string `json:"recipe"`
    // Slug part of the URL
	Slug   string `json:"slug"`

	// full taco api doubles up url stuff for some reaosn
    CommonTaco
    // URL to the baselayer for the taco
	BaseLayerURL string   `json:"base_layer_url"`

    // URL to the mixin of the taco
	MixinURL     string   `json:"mixin_url"`

    // URL to the condiment of the taco
	CondimentURL string   `json:"condiment_url"`

    // URL to the seasoning of the taco
	SeasoningURL string   `json:"seasoning_url"`

    // URL to the shell of the taco
	ShellURL     string   `json:"shell_url"`
}

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
	// incosistent
	if t.shell.Name {
		_, desc = desc + " in " + t.shell.Name + "."
	}

	return desc
}

// Unmarshall a FullTaco to json []byte
func (t *FullTaco) UnmarshalJSON() data []byte, error {
    data []byte
	doc := FullTacoJSON{
        Name: t.Name, URL: t.URL, Recipe: t.Recipe, Slug: t.Slug,
        BaseLayer: t.baseLayer, Mixin: t.mixin, Condiment: t.condiment,
        Seasoning: t.seasoning, Shell: t.shell}
	err := json.Unmarshal(data, doc)
	if err != nil {
		return err
	}

	return data
}

// Provides a struct for marshaling and unmarshaling a FullTaco
type FullTacoJSON struct {
    Name   string `json:"name"`
    // The URL to the recipe
	URL    string `json:"url"`
    // Recipe instructions on how to make this Taco
	Recipe string `json:"recipe"`
    // Slug part of the URL
	Slug   string `json:"slug"`

    RandomTacoJSON
    // full taco response doubles up on these URLs on data from
    // the RandomTaco
    // URL to the baselayer for the taco
	BaseLayerURL string   `json:"base_layer_url"`
    // URL to the mixin of the taco
	MixinURL     string   `json:"mixin_url"`
    // URL to the condiment of the taco
	CondimentURL string   `json:"condiment_url"`
    // URL to the seasoning of the taco
	SeasoningURL string   `json:"seasoning_url"`
    // URL to the shell of the taco
	ShellURL     string   `json:"shell_url"`
}

// Marshal's to a FullTaco
func (t *FullTacoJSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}
