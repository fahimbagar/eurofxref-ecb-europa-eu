package interfaces

import (
	"encoding/json"
	"net/http"
)

type ExchangeAgent interface {
	GetLatestExchange() ForexResponse
	GetExchangeByDate(date string) ForexResponse
	GetAnalyzedRates() ForexResponse
}

type WebserviceHandler struct {
	ExchangeAgent ExchangeAgent
}

type HelloWorld struct {
	Hello string `json:"hello"`
}

func (handler WebserviceHandler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(HelloWorld{Hello: "world"})
}

func (handler WebserviceHandler) GetLatestExchange(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(handler.ExchangeAgent.GetLatestExchange())
}

func (handler WebserviceHandler) GetLatestExchangeByDate(w http.ResponseWriter, r *http.Request) {
	dateString := r.Context().Value("match")
	if dateString == nil {
		http.NotFound(w, r)
	}
	_ = json.NewEncoder(w).Encode(handler.ExchangeAgent.GetExchangeByDate(dateString.(string)))
}

func (handler WebserviceHandler) RatesAnalyze(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(handler.ExchangeAgent.GetAnalyzedRates())
}
