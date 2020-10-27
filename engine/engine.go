package engine

import (
	"encoding/xml"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/domain"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/interfaces"
	"log"
	"net/http"
)

type ExchangeEngine struct {
	BaseCurrency       string
	ExchangeRepository ExchangeRepository
}

type ExchangeRepository interface {
	ResetDB() error
	Store(envelope interfaces.Envelope) error
	FindByLatestDate() []domain.Exchange
	FindByDateString(date string) []domain.Exchange
	Find() []domain.Exchange
}

// ExchangeEngine downloading rates data from ECB and preparing the database
func (exchangeEngine ExchangeEngine) Prepare() {
	client := http.Client{}
	resp, err := client.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml")
	if err != nil {
		log.Fatal(err)
	}

	var envelope interfaces.Envelope
	if err = xml.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		log.Fatal(err)
	}

	log.Println("resetting database")
	if err := exchangeEngine.ExchangeRepository.ResetDB(); err != nil {
		log.Fatal(err)
	}

	if err := exchangeEngine.ExchangeRepository.Store(envelope); err != nil {
		log.Fatal(err)
	}
	log.Println("finish preparing exchange engine (use case)...")
}

// GetLatestExchange returns rates with DTO ForexResponse
func (exchangeEngine *ExchangeEngine) GetLatestExchange() (fxResponse interfaces.ForexResponse) {
	result := exchangeEngine.ExchangeRepository.FindByLatestDate()

	fxResponse.Base = exchangeEngine.BaseCurrency
	fxResponse.Rates = make(map[string]float64, 0)
	for _, fx := range result {
		fxResponse.Rates[fx.Currency] = fx.Rate
		fxResponse.Date = fx.ForexDate.Format("2006-01-02")
	}

	return
}

// GetExchangeByDate returns rates with certain date
func (exchangeEngine ExchangeEngine) GetExchangeByDate(date string) (fxResponse interfaces.ForexResponse) {
	result := exchangeEngine.ExchangeRepository.FindByDateString(date)

	fxResponse.Base = exchangeEngine.BaseCurrency
	fxResponse.Rates = make(map[string]float64, 0)
	for _, fx := range result {
		fxResponse.Rates[fx.Currency] = fx.Rate
		fxResponse.Date = fx.ForexDate.Format("2006-01-02")
	}

	return
}

// GetAnalyzedRates returns analyzed rate: max, min, average
func (exchangeEngine ExchangeEngine) GetAnalyzedRates() (fxResponse interfaces.ForexResponse) {
	result := exchangeEngine.ExchangeRepository.Find()

	fxResponse.Base = exchangeEngine.BaseCurrency
	fxResponse.RatesAnalyzer = make(map[string]interfaces.RatesAnalyze, 0)
	for _, fx := range result {
		rates, ok := fxResponse.RatesAnalyzer[fx.Currency]
		if !ok {
			fxResponse.RatesAnalyzer[fx.Currency] = interfaces.RatesAnalyze{
				Min:   fx.Rate,
				Max:   fx.Rate,
				Sum:   fx.Rate,
				Count: 1,
			}
		} else {
			rates := fxResponse.RatesAnalyzer[fx.Currency]
			rates.Sum = rates.Sum + fx.Rate
			rates.Count = rates.Count + 1
			fxResponse.RatesAnalyzer[fx.Currency] = rates
		}

		if fx.Rate < rates.Min {
			rates := fxResponse.RatesAnalyzer[fx.Currency]
			rates.Min = fx.Rate
			fxResponse.RatesAnalyzer[fx.Currency] = rates
		}

		if fx.Rate > rates.Max {
			rates := fxResponse.RatesAnalyzer[fx.Currency]
			rates.Max = fx.Rate
			fxResponse.RatesAnalyzer[fx.Currency] = rates
		}
	}

	for currency, _ := range fxResponse.RatesAnalyzer {
		rates := fxResponse.RatesAnalyzer[currency]
		rates.Avg = rates.Sum / float64(rates.Count)
		fxResponse.RatesAnalyzer[currency] = rates
	}

	return
}
