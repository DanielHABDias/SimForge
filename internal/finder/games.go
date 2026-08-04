package finder

var KnownGames = map[string]string{

	"1238060": "Dead Space 3",

	"1426210": "It Takes Two",

	"1222680": "Need For Speed Heat",

	"1262560": "Need For Speed Most Wanted",

	"1262580": "Need For Speed Payback",

	"1846380": "Need For Speed Unbound",

	"1172380": "STAR WARS Jedi Fallen Order",

	"1774580": "STAR WARS Jedi Survivor",

	"1222670": "The Sims 4",
}


func GetGameName(appID string) string {

	if name, exists := KnownGames[appID]; exists {
		return name
	}


	return "Unknown prefix (" + appID + ")"
}