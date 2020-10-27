package main

import (
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/engine"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/infrastructure"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/interfaces"
	"log"
	"net/http"
	"regexp"
	"time"
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
	customHandler.HandleFunc(regexp.MustCompile(`/rates/analyze$`), webservices.RatesAnalyze)

	log.Println("list of endpoints:")
	log.Println("- http://localhost:8282/rates/latest")
	log.Printf("- http://localhost:8282/rates/%s | date format is yyyy-mm-dd\n", time.Now().AddDate(0, 0, -2).Format("2006-01-02"))
	log.Println("- http://localhost:8282/rates/analyze")
	log.Println()

	log.Println("web server started at http://0.0.0.0:8282")
	if err = http.ListenAndServe(":8282", customHandler); err != nil {
		log.Fatal(err)
	}
}
