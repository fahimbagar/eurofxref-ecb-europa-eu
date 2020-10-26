package interfaces

import (
	"encoding/xml"
)

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Gesmes  string   `xml:"gesmes,attr"`
	Xmlns   string   `xml:"xmlns,attr"`
	Script  string   `xml:"script"`
	Subject string   `xml:"subject"`
	Sender  struct {
		Text string `xml:",chardata"`
		Name string `xml:"name"`
	} `xml:"Sender"`
	Exchanges struct {
		Text              string `xml:",chardata"`
		CurrenciesPerDate []struct {
			Text     string `xml:",chardata"`
			Time     string `xml:"time,attr"`
			Currency []struct {
				Text     string `xml:",chardata"`
				Currency string `xml:"currency,attr"`
				Rate     string `xml:"rate,attr"`
			} `xml:"Cube"`
		} `xml:"Cube"`
	} `xml:"Cube"`
}
