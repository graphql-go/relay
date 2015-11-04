package starwars

import "fmt"

/*
This defines a basic set of data for our Star Wars Schema.

This data is hard coded for the sake of the demo, but you could imagine
fetching this data from a backend service rather than from hardcoded
JSON objects in a more complex demo.
*/

type Ship struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Faction struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Ships []string `json:"ships"`
}

var XWing = &Ship{"1", "X-Wing"}
var YWing = &Ship{"2", "Y-Wing"}
var AWing = &Ship{"3", "A-Wing"}

// Yeah, technically it's Corellian. But it flew in the service of the rebels,
// so for the purposes of this demo it's a rebel ship.
var Falcon = &Ship{"4", "Millenium Falcon"}
var HomeOne = &Ship{"5", "Home One"}
var TIEFighter = &Ship{"6", "TIE Fighter"}
var TIEInterceptor = &Ship{"7", "TIE Interceptor"}
var Executor = &Ship{"8", "Executor"}

var Rebels = &Faction{
	"1",
	"Alliance to Restore the Republic",
	[]string{"1", "2", "3", "4", "5"},
}

var Empire = &Faction{
	"2",
	"Galactic Empire",
	[]string{"6", "7", "8"},
}

var factions = map[string]*Faction{
	"1": Rebels,
	"2": Empire,
}
var ships = map[string]*Ship{
	"1": XWing,
	"2": YWing,
	"3": AWing,
	"4": Falcon,
	"5": HomeOne,
	"6": TIEFighter,
	"7": TIEInterceptor,
	"8": Executor,
}
var nextShip = 9

func CreateShip(shipName string, factionId string) *Ship {
	nextShip = nextShip + 1
	newShip := &Ship{
		fmt.Sprintf("%v", nextShip),
		shipName,
	}
	ships[newShip.ID] = newShip

	faction := GetFaction(factionId)
	if faction != nil {
		faction.Ships = append(faction.Ships, newShip.ID)
	}
	return newShip
}
func GetShip(id string) *Ship {
	if ship, ok := ships[id]; ok {
		return ship
	}
	return nil
}
func GetFaction(id string) *Faction {
	if faction, ok := factions[id]; ok {
		return faction
	}
	return nil
}
func GetRebels() *Faction {
	return Rebels
}
func GetEmpire() *Faction {
	return Empire
}
