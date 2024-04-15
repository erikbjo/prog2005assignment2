package responses

type ResponseFromRestcountries struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	Cca2       string                `json:"cca2"`
	Currencies map[string]Currencies `json:"currencies"`
	Capital    []string              `json:"capital"`
	Latlng     []float64             `json:"latlng"`
	Area       float64               `json:"area"`
	Population int                   `json:"population"`
}

type Currencies struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}
