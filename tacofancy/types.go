
import "fmt"
// see https://taco-randomizer.herokuapp.com/random/ for a random example
// Seems to be consistent in one base layer one shell, one seasoning,
// one mixin, one condiment
type TacoPart struct {
    Name string
    URL string
    Slug string
    Recipe string
}

type Taco struct {
    BaseLayer TacoPart
    Mixin TacoPart
    Condiment TacoPart
    Seasoning TacoPart
    Shell TacoPart
}

func (t *Taco) Description() string {
    return fmt.Printf(t.BaseLayer)
}
