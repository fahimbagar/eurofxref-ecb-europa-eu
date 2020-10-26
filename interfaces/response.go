package interfaces

type ForexResponse struct {
	Base          string                  `json:"base"`
	Rates         map[string]float64      `json:"rates,omitempty"`
	RatesAnalyzer map[string]RatesAnalyze `json:"rates_analyze,omitempty"`
	Date          string                  `json:"date,omitempty"`
}

type RatesAnalyze struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
	Avg float64 `json:"avg"`
}
