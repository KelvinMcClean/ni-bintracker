package bintracker

type Bin struct {
	Date          string
	DayOfWeekName string
	ColourLabel   string
	Name          string
	Colour        string
	Capacity      string
	BinType       string `json:"type"`
	UPRN          string
}
