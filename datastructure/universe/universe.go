package universe

type Moon struct {
	Text string `xml:",chardata"`
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
	Size string `xml:"size,attr"`
}
type Planet struct {
	Text   string `xml:",chardata"`
	ID     string `xml:"id,attr"`
	Player string `xml:"player,attr"`
	Name   string `xml:"name,attr"`
	Coords string `xml:"coords,attr"`
	Moon   Moon   `xml:"moon"`
}

type Universe struct {
	Planet []Planet `xml:"planet"`
}
