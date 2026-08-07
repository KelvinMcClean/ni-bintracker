package gateway

type Response struct {
	Status   string   `json:"status"`
	Message  string   `json:"message"`
	LastWeek []BinDay `json:"lastWeek"`
	ThisWeek []BinDay `json:"thisWeek"`
	NextWeek []BinDay `json:"nextWeek"`
}

type BinDay struct {
	Date          string `json:"date"`
	DayOfWeekName string `json:"dayOfWeekName"`
	Bins          []Bin  `json:"bins"`
}

type Bin struct {
	ColourLabel string `json:"colourLabel"`
	Colour      string `json:"colour"`
	Name        string `json:"name"`
	Capacity    string `json:"capacity"`
	BinType     string `json:"type"`
	UPRN        string `json:"uprn"`
}
