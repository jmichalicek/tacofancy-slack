package tacofancy


// An interface implementing a tco
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
type RandomTaco struct {
    Name   string `json:"name"`
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

// Returns a description of the taco made up from the names of the parts
func (t *RandomTaco) Description() string {

	// shell names are inconsistent, but roll with this for now
	desc := t.BaseLayer.Name+" seasoned with "+t.Seasoning.Name+" with "+t.Mixin.Name+" and "+
		t.Condiment.Name+" in "+t.Shell.Name+"."
	return desc
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

    // URL to the baselayer for the taco
	BaseLayerURL string   `json:"base_layer_url"`
    // The actual baselayer
	BaseLayer    TacoPart `json:"base_layer"`

    // URL to the mixin of the taco
	MixinURL     string   `json:"mixin_url"`
    // The mixin
	Mixin        TacoPart `json:"mixin"`

    // URL to the condiment of the taco
	CondimentURL string   `json:"condiment_url"`
    // The condiment
	Condiment    TacoPart `json:"condiment"`

    // URL to the seasoning of the taco
	SeasoningURL string   `json:"seasoning_url"`
    // The seasoning
	Seasoning    TacoPart `json:"seasoning"`

    // URL to the shell of the taco
	ShellURL     string   `json:"shell_url"`
    // The shell
	Shell        TacoPart `json:"shell"`
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
	 if t.Shell.Name {
	 	desc = desc + " in " + t.Shell.Name + "."
	 }

	return desc
}
