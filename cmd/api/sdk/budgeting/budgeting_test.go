package budgeting_test

import (
	"testing"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elopez00/scale-backend/pkg/test"
	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/elopez00/scale-backend/cmd/api/middleware"
	"github.com/elopez00/scale-backend/cmd/api/sdk/budgeting"
)

var user = models.User {
	Id: "testvalue",
}

var budget = models.Budget{
	Categories: []models.Category {
		{ Name: "shopping", Budget: 200 },
		{ Name: "groceries", Budget: 250 },
		{ Name: "rent", Budget: 800 },
	},

	WhiteList: []models.WhiteListItem {
		{ Category: "shopping", Name: "Calvin Klien" },
		{ Category: "shopping", Name: "Best Buy" },
		{ Category: "shopping", Name: "Amazon" },
		{ Category: "groceries", Name: "Aldi" },
		{ Category: "groceries", Name: "Walmart" },
		{ Category: "rent", Name: "The Rise" },
	},

	Request: models.UpdateRequest {
		Update: models.UpdateObject {
			Categories: []models.Category {
				{ Name: "shopping", Budget: 200, Id: "qwert" },
				{ Name: "groceries", Budget: 250, Id: "asdfag" },
				{ Name: "rent", Budget: 800, Id: ";lkjk" },
			},
			WhiteList: []models.WhiteListItem {
				{ Category: "shopping", Name: "Calvin Klien", Id: ";lkjl" },
				{ Category: "shopping", Name: "Best Buy", Id: "asdfasdf" },
				{ Category: "shopping", Name: "Amazon", Id: "qwerqwer" },
				{ Category: "groceries", Name: "Aldi", Id: ";sdlfkgjsd" },
				{ Category: "groceries", Name: "Walmart", Id: "zxcvzxc" },
				{ Category: "rent", Name: "The Rise", Id: ".,mn.,n,mn" },
			},
		},
	},
}

func TestCreateBudget(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close()

	jsonObject, _ := json.Marshal(budget)

	// test categories query
	query1 := 
		`INSERT INTO categories\(id, name, budget, categoryId\) ` +
		`VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) AS updated ON DUPLICATE KEY UPDATE ` +
		`id\=updated\.id, name\=updated\.name, budget\=updated\.budget, categoryId\=updated\.categoryId;`
	mock.ExpectPrepare(query1).
	ExpectExec().
	WillReturnResult(sqlmock.NewResult(0, 0))
	
	// test whitelist query
	query2 := 
		`INSERT INTO whitelist\(id, category, name, itemId\) ` +
		`VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\), \(\?,\?,\?,\?\), \(\?,\?,\?,\?\), \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) AS updated ON DUPLICATE KEY UPDATE ` +
		`id\=updated\.id, category\=updated\.category, name\=updated\.name, itemId\=updated\.itemId;`
	mock.ExpectPrepare(query2).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(0, 0))
	
	res := test.PostWithCookie(
		"/v0/createBudget",
		middleware.Authenticate(budgeting.Create(app), app),
		bytes.NewBuffer(jsonObject),
		app,
		"AuthToken",
	)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error("There were unfulfilled expectations:", err)
		return
	}

	if res.Code != http.StatusOK {
		// body response
		var response models.Response 
		json.NewDecoder(res.Body).Decode(&response)

		t.Errorf("Expected status to be %v, instead we got %v with error message: %v", http.StatusOK, res.Code, response.Message)
		return
	}
}