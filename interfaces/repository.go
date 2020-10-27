package interfaces

import (
	"fmt"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/domain"
	"log"
)

type DBHandler interface {
	Execute(statement string) error
	Query(statement string) Row
}

type Row interface {
	Scan(dest ...interface{}) error
	Next() bool
}

type DBRepository struct {
	dbHandlers map[string]DBHandler
	dbHandler  DBHandler
}

type ExchangeRepository DBRepository

// NewDBExchange init database handler
func NewDBExchange(dbHandlers map[string]DBHandler) (*ExchangeRepository, error) {
	dbExchange := new(ExchangeRepository)
	dbExchange.dbHandlers = dbHandlers
	dbExchange.dbHandler = dbHandlers["ExchangeRepository"]
	return dbExchange, nil
}

// ResetDB resetting the database
func (repo *ExchangeRepository) ResetDB() error {
	if err := repo.dbHandler.Execute("DROP TABLE IF EXISTS exchange;"); err != nil {
		return err
	}

	if err := repo.dbHandler.Execute(`
		CREATE TABLE exchange
		(   
			id INTEGER CONSTRAINT exchange_pk PRIMARY KEY AUTOINCREMENT,
			currency TEXT NOT NULL,
			forex_date TIMESTAMP NOT NULL,
			rate FLOAT,
			createdAt DATE DEFAULT CURRENT_TIMESTAMP NOT NULL
		);
	`); err != nil {
		return err
	}

	return nil
}

// Store data from ECB rates to database
func (repo *ExchangeRepository) Store(envelope Envelope) error {
	for _, currencies := range envelope.Exchanges.CurrenciesPerDate {
		for _, currency := range currencies.Currency {
			if err := repo.dbHandler.Execute(fmt.Sprintf("INSERT INTO exchange (currency, rate, forex_date) VALUES ('%s', '%s', '%s')", currency.Currency, currency.Rate, currencies.Time)); err != nil {
				return err
			}
		}
	}

	return nil
}

// FindByLatestDate returns exchanges with latest date from database
func (repo *ExchangeRepository) FindByLatestDate() []domain.Exchange {
	row := repo.dbHandler.Query(`
		SELECT t.currency, t.rate, t.forex_date
		FROM exchange t
				 INNER JOIN (
			SELECT currency, max(forex_date) AS MaxDate
			FROM exchange
			GROUP BY currency
		) tm ON t.currency = tm.currency AND t.forex_date = tm.MaxDate
		ORDER BY t.currency;
	`)
	var exchanges []domain.Exchange
	for row.Next() {
		var exchange domain.Exchange
		if err := row.Scan(&exchange.Currency, &exchange.Rate, &exchange.ForexDate); err != nil {
			log.Fatal(err)
		}
		exchanges = append(exchanges, exchange)
	}

	return exchanges
}

// FindByDateString returns exchanges
func (repo *ExchangeRepository) FindByDateString(date string) []domain.Exchange {
	row := repo.dbHandler.Query(fmt.Sprintf(`
		SELECT t.currency, t.rate, t.forex_date
		FROM exchange t
		WHERE t.forex_date = '%s'
		ORDER BY t.currency;
	`, date))
	var exchanges []domain.Exchange
	for row.Next() {
		var exchange domain.Exchange
		if err := row.Scan(&exchange.Currency, &exchange.Rate, &exchange.ForexDate); err != nil {
			log.Fatal(err)
		}
		exchanges = append(exchanges, exchange)
	}

	return exchanges
}

// Find returns all exchanges
func (repo *ExchangeRepository) Find() []domain.Exchange {
	row := repo.dbHandler.Query(`
		SELECT currency, rate, forex_date
		FROM exchange
		ORDER BY currency;
	`)
	var exchanges []domain.Exchange
	for row.Next() {
		var exchange domain.Exchange
		if err := row.Scan(&exchange.Currency, &exchange.Rate, &exchange.ForexDate); err != nil {
			log.Fatal(err)
		}
		exchanges = append(exchanges, exchange)
	}

	return exchanges
}
