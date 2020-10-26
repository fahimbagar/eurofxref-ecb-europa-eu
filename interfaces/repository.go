package interfaces

import (
	"fmt"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/domain"
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

type IExchange DBRepository

func NewDBExchange(dbHandlers map[string]DBHandler) (*IExchange, error) {
	dbExchange := new(IExchange)
	dbExchange.dbHandlers = dbHandlers
	dbExchange.dbHandler = dbHandlers["ExchangeRepository"]
	return dbExchange, nil
}

func (repo *IExchange) ResetDB() error {
	if err := repo.dbHandler.Execute("DROP TABLE IF EXISTS exchange;"); err != nil {
		return err
	}

	if err := repo.dbHandler.Execute(`
		CREATE TABLE exchange
		(   
			id INTEGER CONSTRAINT exchange_pk PRIMARY KEY AUTOINCREMENT,
			currency TEXT NOT NULL,
			forex_date DATE NOT NULL,
			rate FLOAT,
			createdAt DATE DEFAULT CURRENT_TIMESTAMP NOT NULL
		);
	`); err != nil {
		return err
	}

	return nil
}

func (repo *IExchange) Store(envelope domain.Envelope) error {
	for _, currencies := range envelope.Exchanges.CurrenciesPerDate {
		for _, currency := range currencies.Currency {
			if err := repo.dbHandler.Execute(fmt.Sprintf("INSERT INTO exchange (currency, forex_date, rate) VALUES ('%s', '%s', '%s')", currency.Currency, currencies.Time, currency.Rate)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (repo *IExchange) FindByLatestDate() []domain.Exchange {
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
			panic(err)
		}
		exchanges = append(exchanges, exchange)
	}

	return exchanges
}

func (repo *IExchange) FindByDateString(date string) []domain.Exchange {
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
			panic(err)
		}
		exchanges = append(exchanges, exchange)
	}

	return exchanges
}
