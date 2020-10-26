package engine

import (
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/infrastructure"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/interfaces"
	"log"
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
	log.Println(testDate)

	forexResponse = agent.GetExchangeByDate(testDate)
	if len(forexResponse.Rates) == 0 {
		t.Error()
	}

	forexResponse = agent.GetExchangeByDate("2010-01-01")
	if len(forexResponse.Rates) > 0 {
		t.Error()
	}
}