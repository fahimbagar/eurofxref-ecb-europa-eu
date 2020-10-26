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
}

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

func (exchangeEngine *ExchangeEngine) GetLatestExchange() interfaces.ForexResponse {
	forex := exchangeEngine.ExchangeRepository.FindByLatestDate()

	exchangesResponse := interfaces.ForexResponse{
		Base:  exchangeEngine.BaseCurrency,
		Rates: make(map[string]float64, 0),
	}
	for _, fx := range forex {
		exchangesResponse.Rates[fx.Currency] = fx.Rate
		exchangesResponse.Date = fx.ForexDate.Format("2006-01-02")
	}

	return exchangesResponse
}

func (exchangeEngine ExchangeEngine) GetExchangeByDate(date string) interfaces.ForexResponse {
	forex := exchangeEngine.ExchangeRepository.FindByDateString(date)

	exchangesResponse := interfaces.ForexResponse{
		Base:  exchangeEngine.BaseCurrency,
		Rates: make(map[string]float64, 0),
	}
	for _, fx := range forex {
		exchangesResponse.Rates[fx.Currency] = fx.Rate
		exchangesResponse.Date = fx.ForexDate.Format("2006-01-02")
	}

	return exchangesResponse
}
