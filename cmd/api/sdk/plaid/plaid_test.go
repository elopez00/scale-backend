package plaid_test

import (
	"encoding/json"
	"net/http"
	"testing"

	m "github.com/elopez00/scale-backend/cmd/api/middleware"
	p "github.com/elopez00/scale-backend/cmd/api/sdk/plaid"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/elopez00/scale-backend/pkg/test"
)

var token = models.Token {
	Value: 	"access-sandbox-3b6a6577-4c02-4fc3-a213-b8adf828c38f",
	Id: 	"nothin",
}

var user = models.User {
	Id: "testvalue",
}

func TestInvalidClient(t *testing.T) {
	app, _ := test.GetMockApp() 
	defer app.DB.Client.Close()


	if res := test.GetWithCookie(
		"/v0/getLinkToken",
		m.Authenticate(p.GetPlaidToken(app), app), 
		nil, 
		app, 
		"AuthToken",
	); res.Code != http.StatusBadGateway {
		t.Errorf("Failed get. Expected %v, instead got %v", http.StatusOK, res.Code)
	} else {
		var response models.Response
		json.NewDecoder(res.Body).Decode(&response)

		if response.Message != "Failure to load client" {
			t.Errorf("Link token shouldn't have been extracted, instead recieved error: %v", response.Message)
		}
	}
}

func TestLinkTokenRetrieval(t *testing.T) {
	app, _ := test.GetPlaidMockApp()
	defer app.DB.Client.Close()

	if res := test.GetWithCookie(
		"/v0/getLinkToken",
		m.Authenticate(p.GetPlaidToken(app), app), 
		nil, 
		app, 
		"AuthToken",
	); res.Code != http.StatusOK {
		t.Errorf("Failed get. Expected %v, instead got %v", http.StatusOK, res.Code)
	} else {
		var response models.Response
		json.NewDecoder(res.Body).Decode(&response)

		if response.Message != "Successfully recieved link token from plaid" {
			t.Errorf("Link token was not extracted successfuly, instead recieved error: %v", response.Message)
		}
	}
}

func TestGetTransactions(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	rows := sqlmock.NewRows([]string{"id", "token", "itemID"}).
		AddRow(user.Id, token.Value, token.Id)
	
	query := `SELECT id, token, itemID FROM plaidtokens WHERE id\="testvalue"`
	mock.ExpectQuery(query).WillReturnRows(rows)

	res := test.GetWithCookie(
		"/v0/getTransactions",
		m.Authenticate(p.GetTransactions(app), app),
		nil,
		app,
		"AuthToken",
	)

	if res.Code != http.StatusOK {
		var response models.Response
		json.NewDecoder(res.Body).Decode(&response)

		t.Errorf("This call returned the wrong http status. Expected %v, got %v", http.StatusOK, res.Code)
		t.Error("The call did not return the intended result, instead", response.Message)
	}
}