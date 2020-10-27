package interfaces_test

import (
	"encoding/xml"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/infrastructure"
	"github.com/fahimbagar/eurofxref-ecb-europa-eu/interfaces"
	"io/ioutil"
	"os"
	"testing"
)

func Test_Repository(t *testing.T) {
	dbHandler := infrastructure.NewSqliteHandler("test.db")

	handlers := make(map[string]interfaces.DBHandler)
	handlers["ExchangeRepository"] = dbHandler

	dbConn, err := interfaces.NewDBExchange(handlers)
	if err != nil {
		t.Error(err)
	}

	if err := dbConn.ResetDB(); err != nil {
		t.Error(err)
	}

	xmlFile, err := os.Open("response_test.xml")
	if err != nil {
		t.Error(err)
	}
	var envelope interfaces.Envelope
	byteValue, _ := ioutil.ReadAll(xmlFile)
	if err := xml.Unmarshal(byteValue, &envelope); err != nil {
		t.Error(err)
	}

	if err := xmlFile.Close(); err != nil {
		t.Error(err)
	}

	if err := dbConn.Store(envelope); err != nil {
		t.Error(err)
	}

	exchanges := dbConn.FindByLatestDate()
	if len(exchanges) == 0 {
		t.Error()
	}

	exchanges = dbConn.FindByDateString("2020-10-30")
	if len(exchanges) > 0 {
		t.Error()
	}

	exchanges = dbConn.FindByDateString("2020-10-22")
	if len(exchanges) == 0 {
		t.Error()
	}

	exchanges = dbConn.Find()
	if len(exchanges) == 0 {
		t.Error()
	}

	if err := os.Remove("test.db"); err != nil {
		t.Error(err)
	}
}
