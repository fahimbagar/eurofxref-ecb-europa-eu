package main

import (
	"encoding/xml"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/domain"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/engine"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/infrastructure"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/interfaces"
	"log"
	"net/http"
	"regexp"
)

func main() {
	log.Print("starting...")
	client := http.Client{}
	resp, err := client.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml")
	if err != nil {
		log.Fatal(err)
	}

	var envelope domain.Envelope
	if err = xml.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		log.Fatal(err)
	}

	dbHandler := infrastructure.NewSqliteHandler("exchange.db")
	log.Print("creating db...")

	handlers := make(map[string]interfaces.DBHandler)
	handlers["ExchangeRepository"] = dbHandler

	agent := new(engine.ExchangeAgent)
	agent.BaseCurrency = "EUR"
	agent.ExchangeRepository, err = interfaces.NewDBExchange(handlers)
	if err != nil {
		log.Fatal(err)
	}

	agent.Reset()
	agent.Store(envelope)

	customHandler := &interfaces.Middleware{}

	webservices := interfaces.WebserviceHandler{}
	webservices.ExchangeAgent = agent

	customHandler.HandleFunc(regexp.MustCompile(`/hello-world$`), webservices.HelloWorld)
	customHandler.HandleFunc(regexp.MustCompile(`/rates/latest$`), webservices.GetLatestExchange)
	customHandler.HandleFunc(regexp.MustCompile(`/rates/(\d{4}-\d{2}-\d{2})$`), webservices.GetLatestExchangeByDate)

	log.Println("http served at http://localhost:8080")
	if err = http.ListenAndServe(":8080", customHandler); err != nil {
		log.Fatal(err)
	}

	//for _, currencies := range envelope.Exchanges.CurrenciesPerDate {
	//	for _, currency := range currencies.Currency {
	//		fmt.Printf("INSERT INTO exchange (currency, rate, forex_date) VALUES ('%s', '%s', '%s');\n", currency.Currency, currency.Rate, currencies.Time)
	//	}
	//}
	//
	//forex := engine.ForexResponse{
	//	Base:          "",
	//	Rates:         make(map[string]float64),
	//}
	//for _, currency := range envelope.Exchanges.CurrenciesPerDate[0].Currency {
	//	forex.Rates[currency.Currency], err = strconv.ParseFloat(currency.Rate, 64)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
	//
	//p, err := json.MarshalIndent(&forex, "", "   ")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("%s", p)
}
