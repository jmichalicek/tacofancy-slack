package tacofancy

//import "fmt"
// see https://taco-randomizer.herokuapp.com/random/ for a random example
// Seems to be consistent in one base layer one shell, one seasoning,
// one mixin, one condiment
type TacoPart struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Slug   string `json:"slug"`
	Recipe string `json:"recipe"`
}

type RandomTaco struct {
	BaseLayer TacoPart `json:"base_layer"`
	Mixin     TacoPart `json:"mixin"`
	Condiment TacoPart `json:"condiment"`
	Seasoning TacoPart `json:"seasoning"`
	Shell     TacoPart `json:"shell"`
}

func (t *RandomTaco) Description() string {

	// shell names are inconsistent, but roll with this for now
	_, desc := t.BaseLayer.Name+" seasoned with ", t.Seasoning.Name+" with "+t.Mixin.Name+" and "+
		t.Condiment.Name+" in "+t.Shell.Name+"."
	return desc
}

type FullTaco struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Recipe string `json:"recipe"`
	Slug   string `json:"slug"`

	// full taco api doubles up url stuff for some reaosn
	BaseLayerURL string   `json:"base_layer_url"`
	BaseLayer    TacoPart `json:"base_layer"`
	MixinURL     string   `json:"mixin_url"`
	Mixin        TacoPart `json:"mixin"`
	CondimentURL string   `json:"condiment_url"`
	Condiment    TacoPart `json:"condiment"`
	SeasoningURL string   `json:"seasoning_url"`
	Seasoning    TacoPart `json:"seasoning"`
	ShellURL     string   `json:"shell_url"`
	Shell        TacoPart `json:"shell"`
}

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
	// incosistent
	// if t.Shell.Name {
	// 	_, desc = desc + " in " + t.Shell.Name + "."
	// }

	return desc
}
