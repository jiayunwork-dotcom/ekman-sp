package config

type ExampleCase struct {
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	LatitudeDeg float64 `json:"latitude_deg"`
	WindDirDeg  float64 `json:"wind_dir_deg"`
	Tau         float64 `json:"tau"`
	Description string  `json:"description"`
}

var Catalog = []ExampleCase{
	{
		Name:        "midlat-wind",
		Path:        "example/midlat-wind.json",
		LatitudeDeg: 45,
		WindDirDeg:  90,
		Tau:         0.2,
		Description: "45N mid-latitude westerly wind",
	},
	{
		Name:        "south-hemi-wind",
		Path:        "example/south-hemi-wind.json",
		LatitudeDeg: -40,
		WindDirDeg:  270,
		Tau:         0.15,
		Description: "40S southern hemisphere westerly",
	},
	{
		Name:        "north-trade",
		Path:        "example/north-trade.json",
		LatitudeDeg: 15,
		WindDirDeg:  270,
		Tau:         0.12,
		Description: "15N tropical easterlies",
	},
}

func CatalogJSON() []ExampleCase {
	out := make([]ExampleCase, len(Catalog))
	copy(out, Catalog)
	return out
}

func FindExample(name string) (ExampleCase, bool) {
	for _, c := range Catalog {
		if c.Name == name {
			return c, true
		}
	}
	return ExampleCase{}, false
}

func LoadExample(name string) (*Resolved, error) {
	c, ok := FindExample(name)
	if !ok {
		return nil, ErrUnknownExample
	}
	return Load(c.Path)
}
