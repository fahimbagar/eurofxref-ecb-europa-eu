package main

import (
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/engine"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/infrastructure"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/interfaces"
	"log"
	"net/http"
	"regexp"
)

func main() {
	log.Print("starting...")

	dbHandler := infrastructure.NewSqliteHandler("exchange.db")
	log.Print("creating db...")

	handlers := make(map[string]interfaces.DBHandler)
	handlers["ExchangeRepository"] = dbHandler

	agent := new(engine.ExchangeEngine)
	agent.BaseCurrency = "EUR"
	var err error
	agent.ExchangeRepository, err = interfaces.NewDBExchange(handlers)
	if err != nil {
		log.Fatal(err)
	}
	agent.Prepare()

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
}
