package engine

import (
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/domain"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/interfaces"
	"log"
)

type ExchangeAgent struct {
	BaseCurrency       string
	ExchangeRepository ExchangeRepository
}

type ExchangeRepository interface {
	ResetDB() error
	Store(envelope domain.Envelope) error
	FindByLatestDate() []domain.Exchange
	FindByDateString(date string) []domain.Exchange
}

func (agent ExchangeAgent) Reset() {
	log.Println("resetting database")
	if err := agent.ExchangeRepository.ResetDB(); err != nil {
		log.Fatal(err)
	}
}

func (agent *ExchangeAgent) Store(envelope domain.Envelope) {
	if err := agent.ExchangeRepository.Store(envelope); err != nil {
		log.Fatal(err)
	}
	log.Println("finish preparing agent exchange...")
}

func (agent *ExchangeAgent) GetLatest() interfaces.ForexResponse {
	forex := agent.ExchangeRepository.FindByLatestDate()

	exchangesResponse := interfaces.ForexResponse{
		Base:  agent.BaseCurrency,
		Rates: make(map[string]float64, 0),
	}
	for _, fx := range forex {
		exchangesResponse.Rates[fx.Currency] = fx.Rate
	}

	return exchangesResponse
}

func (agent ExchangeAgent) GetByDate(date string) interfaces.ForexResponse {
	forex := agent.ExchangeRepository.FindByDateString(date)

	exchangesResponse := interfaces.ForexResponse{
		Base:  agent.BaseCurrency,
		Rates: make(map[string]float64, 0),
	}
	for _, fx := range forex {
		exchangesResponse.Rates[fx.Currency] = fx.Rate
	}

	return exchangesResponse
}
