package engine

import (
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/infrastructure"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/interfaces"
	"log"
	"math"
	"testing"
)

func Test_ExchangeEngine(t *testing.T) {
	dbHandler := infrastructure.NewSqliteHandler("test.db")
	handlers := make(map[string]interfaces.DBHandler)
	handlers["ExchangeRepository"] = dbHandler

	agent := new(ExchangeEngine)
	agent.BaseCurrency = "EUR"
	var err error
	agent.ExchangeRepository, err = interfaces.NewDBExchange(handlers)
	if err != nil {
		log.Fatal(err)
	}

	agent.Prepare()

	forexResponse := agent.GetLatestExchange()
	if len(forexResponse.Rates) == 0 {
		t.Error()
	}

	if forexResponse.Base == "" {
		t.Error()
	}

	testDate := forexResponse.Date

	forexResponse = agent.GetExchangeByDate(testDate)
	if len(forexResponse.Rates) == 0 {
		t.Error()
	}

	forexResponse = agent.GetExchangeByDate("2010-01-01")
	if len(forexResponse.Rates) > 0 {
		t.Error()
	}

	forexResponse = agent.GetAnalyzedRates()
	if len(forexResponse.RatesAnalyzer) == 0 {
		t.Error()
	}

	for currency, rates := range forexResponse.RatesAnalyzer {
		if math.Round(rates.Avg * 10000) / 10000 > rates.Max {
			t.Error(currency)
		}

		if math.Round(rates.Avg * 10000) / 10000 < rates.Min {
			t.Error(currency)
		}

		if rates.Min == 0 {
			t.Error(currency)
		}

		if rates.Max == 0 {
			t.Error(currency)
		}

		if rates.Avg == 0 {
			t.Error(currency)
		}
	}
}
