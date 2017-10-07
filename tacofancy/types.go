package tacofancy

// A bit of overengineering going on here just to learn all about
// go's interfaces, embedded structs, etc. and techniques for working
// with them.

// An interface implementing a taco
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
    // To avoid doubling up having to make a private attribute here
    // and public versions on the JSON version we instead just have
    // a private attribute for the JSON
    data randomTacoJSON
}

func (t CommonTaco) BaseLayer() *TacoPart {return &t.data.BaseLayer}
func (t CommonTaco) Mixin() *TacoPart {return &t.data.Mixin}
func (t CommonTaco) Condiment() *TacoPart {return &t.data.Condiment}
func (t CommonTaco) Seasoning() *TacoPart {return &t.data.Seasoning}
func (t CommonTaco) Shell() *TacoPart {return &t.data.Shell}

func (t *CommonTaco) SetBaseLayer(baseLayer *TacoPart)  {t.data.BseLayer = *baseLayer}
func (t *CommonTaco) SetMixin(mixin *TacoPart)  {t.data.Mixin = *mixin}
func (t *CommonTaco) SetCondiment(condiment *TacoPart)  {t.data.Condiment = *condiment}
func (t *CommonTaco) SetSeasoning(seasoning *TacoPart)  {t.data.Seasoning = *seasoning}
func (t *CommonTaco) SetShell(shell *TacoPart)  {t.data.Shell = shell}

// Make a new RandomTaco object
func NewRandomTaco() RandomTaco {
    // pass in the data here
    data = randomTacoJSON{}
    return RandomTaco{data: data}
}

// Implementation of a Taco made of random parts
// This has a BaseLayer, Mixin, Condiment, Seasoning, and Shell
// TacoPart
type RandomTaco struct {
    // this one is marshallable
    CommonTaco
}

// Returns a description of the taco made up from the names of the parts
func (t *RandomTaco) Description() string {
	// shell names are inconsistent, but roll with this for now
	_, desc := t.baseLayer.Name+" seasoned with "+t.seasoning.Name+" with "+t.mixin.Name+" and "+
		t.condiment.Name+" in "+t.shell.Name+"."
	return desc
}

// Implements json.Unmarshaler
// https://golang.org/pkg/encoding/json/#RawMessage.UnmarshalJSON
// https://golang.org/pkg/encoding/json/#Unmarshal calls this method
// Technique taken from https://play.golang.org/p/rQu1W5RXTy linked from
// http://grokbase.com/p/gg/golang-nuts/143vtcxz29/go-nuts-patterns-for-json-encoding-unexported-fields
func (t *RandomTaco) UnmarshalJSON(data []byte) error {
	doc := randomTacoJSON{}
	err := json.Unmarshal(data, &t.doc)
	if err != nil {
		return err
	}

	return nil
}

func (t *RandomTaco) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.data)
}

// Provides a struct for marshaling and unmarshaling a RandomTaco data
// while leaving its fields private so that it is easier to use with Taco interface
type randomTacoJSON struct {
    BaseLayer    TacoPart `json:"base_layer"`
    // The mixin
    Mixin        TacoPart `json:"mixin"`
    Condiment    TacoPart `json:"condiment"`
    // The seasoning
    Seasoning    TacoPart `json:"seasoning"`
    // The shell
    Shell        TacoPart `json:"shell"`
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

// Implements json.Unmarshaler
// https://golang.org/pkg/encoding/json/#RawMessage.UnmarshalJSON
// https://golang.org/pkg/encoding/json/#Unmarshal calls this method
// Technique taken from https://play.golang.org/p/rQu1W5RXTy linked from
// http://grokbase.com/p/gg/golang-nuts/143vtcxz29/go-nuts-patterns-for-json-encoding-unexported-fields
func (t *FullTaco) UnmarshalJSON(data []byte) error {
	doc := fullTacoJSON{}
	err := json.Unmarshal(data, &t.doc)
	if err != nil {
		return err
	}

	return nil
}

func (t *FullTaco) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.data)
}

// Provides a struct for marshaling and unmarshaling a FullTaco
type fullTacoJSON struct {
    Name   string `json:"name"`
    // The URL to the recipe
	URL    string `json:"url"`
    // Recipe instructions on how to make this Taco
	Recipe string `json:"recipe"`
    // Slug part of the URL
	Slug   string `json:"slug"`

    randomTacoJSON
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
