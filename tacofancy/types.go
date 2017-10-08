package tacofancy


// An interface implementing a tco
// I'm not even reall using this
type Taco interface {
    Description() string
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

// Implementation of a Taco made of random parts
// This has a BaseLayer, Mixin, Condiment, Seasoning, and Shell
// TacoPart
type BaseTaco struct {
    // The actual baselayer
	BaseLayer TacoPart `json:"base_layer"`
    // The actual baselayer
	Mixin     TacoPart `json:"mixin"`
    // The actual baselayer
	Condiment TacoPart `json:"condiment"`
    // The actual baselayer
	Seasoning TacoPart `json:"seasoning"`
    // The actual baselayer
	Shell     TacoPart `json:"shell"`
}

// A taco which is not from a recipe, but from random selection of parts
type RandomTaco struct {
    BaseTaco
}

// Returns a description of the taco made up from the names of the parts
func (t *RandomTaco) Description() string {

	// shell names are inconsistent, but roll with this for now
	desc := t.BaseLayer.Name+" seasoned with "+t.Seasoning.Name+" with "+t.Mixin.Name+" and "+
		t.Condiment.Name+" in "+t.Shell.Name+"."
	return desc
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
    BaseTaco
    // The name of the taco
	Name   string `json:"name"`
    // The URL to the recipe
	URL    string `json:"url"`
    // Recipe instructions on how to make this Taco
	Recipe string `json:"recipe"`
    // Slug part of the URL
	Slug   string `json:"slug"`

	// full taco api doubles up url stuff for some reaosn

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
	desc := t.BaseLayer.Name
	if t.Seasoning.Name != "" {
		desc = desc + " seasoned with " + t.Seasoning.Name
	}

	if t.Mixin.Name != "" {
		desc = desc + " with " + t.Mixin.Name
	}

	if t.Condiment.Name != "" {
		joinWord := " with "
		if t.Mixin.Name != "" {
			joinWord = " and "
		}
		desc = desc + joinWord + t.Condiment.Name
	}

	// full tacos do not seem to have a shell usually and shell names are
	// inconsistent
	 if t.Shell.Name != "" {
	 	desc = desc + " in " + t.Shell.Name + "."
	 }

	return desc
}
