package budget_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"fmt"

	"github.com/elopez00/scale-backend/cmd/api/middleware"
	"github.com/elopez00/scale-backend/cmd/api/models"
	budgeting "github.com/elopez00/scale-backend/cmd/api/sdk/budget"
	"github.com/elopez00/scale-backend/pkg/test"

	"github.com/DATA-DOG/go-sqlmock"
)

var user = models.User{
	Id: "testvalue",
}

var budget = models.Budget{
	Categories: []models.Category{
		{Name: "shopping", Budget: 200, WhiteList: []models.WhiteListItem{
			{Category: "shopping", Name: "Calvin Klien"},
			{Category: "shopping", Name: "Best Buy"},
			{Category: "shopping", Name: "Amazon"},
		}},
		{Name: "groceries", Budget: 250, WhiteList: []models.WhiteListItem{
			{Category: "groceries", Name: "Aldi"},
			{Category: "groceries", Name: "Walmart"},
		}},
		{Name: "rent", Budget: 800, WhiteList: []models.WhiteListItem{{Category: "rent", Name: "The Rise"}}},
	},

	Request: models.UpdateRequest{
		Update: models.UpdateObject{
			Categories: []models.Category{
				{Name: "shopping", Budget: 200, Id: "qwert"},
				{Name: "groceries", Budget: 250, Id: "asdfag"},
				{Name: "rent", Budget: 800, Id: ";lkjk"},
			},
			WhiteList: []models.WhiteListItem{
				{Category: "shopping", Name: "Calvin Klien", Id: ";lkjl"},
				{Category: "shopping", Name: "Best Buy", Id: "asdfasdf"},
				{Category: "shopping", Name: "Amazon", Id: "qwerqwer"},
				{Category: "groceries", Name: "Aldi", Id: ";sdlfkgjsd"},
				{Category: "groceries", Name: "Walmart", Id: "zxcvzxc"},
				{Category: "rent", Name: "The Rise", Id: ".,mn.,n,mn"},
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
			`VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) ` +
			`AS updated ON DUPLICATE KEY UPDATE ` +
			`id\=updated\.id, name\=updated\.name, budget\=updated\.budget, categoryId\=updated\.categoryId;`
	mock.
		ExpectPrepare(query1).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(0, 0))

	// test whitelist query
	query2 :=
		`INSERT INTO whitelist\(id, category, name, itemId\) ` +
			`VALUES \(\?,\?,\?,\?\), \(\?,\?,\?,\?\), \(\?,\?,\?,\?\), \(\?,\?,\?,\?\), \(\?,\?,\?,\?\), \(\?,\?,\?,\?\) ` +
			`AS updated ON DUPLICATE KEY UPDATE ` +
			`id\=updated\.id, category\=updated\.category, name\=updated\.name, itemId\=updated\.itemId;`
	mock.
		ExpectPrepare(query2).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(0, 0))

	res := test.PostWithCookie(
		"/v0/createBudget",
		middleware.Authenticate(budgeting.Update(app), app),
		bytes.NewBuffer(jsonObject),
		app,
		"AuthToken",
	)

	test.Response(t, res, http.StatusOK)
	test.MockExpectations(t, mock)
}

func TestGetBudget(t *testing.T) {
	app, mock := test.GetMockApp()
	defer app.DB.Client.Close() 

	categories := budget.Categories
	whitelist := []models.WhiteListItem {
		budget.Categories[0].WhiteList[0],
		budget.Categories[0].WhiteList[1],
		budget.Categories[0].WhiteList[2],
		budget.Categories[1].WhiteList[0],
		budget.Categories[1].WhiteList[1],
		budget.Categories[2].WhiteList[0],
	}

	rows1 := sqlmock.NewRows([]string{"id", "name", "budget", "categoryId"}).
		AddRow(user.Id, categories[0].Name, categories[0].Budget, categories[0].Id).
		AddRow(user.Id, categories[1].Name, categories[1].Budget, categories[1].Id).
		AddRow(user.Id, categories[2].Name, categories[2].Budget, categories[2].Id)

	query1 := fmt.Sprintf("SELECT id, name, budget, categoryId FROM categories WHERE categories.id \\= %q", user.Id)
	mock.
		ExpectQuery(query1).
		WillReturnRows(rows1)
	
	rows2 := sqlmock.NewRows([]string{"id", "name", "category", "itemId"}).
		AddRow(user.Id, whitelist[0].Name, whitelist[0].Category, whitelist[0].Id).
		AddRow(user.Id, whitelist[1].Name, whitelist[1].Category, whitelist[1].Id).
		AddRow(user.Id, whitelist[2].Name, whitelist[2].Category, whitelist[2].Id).
		AddRow(user.Id, whitelist[3].Name, whitelist[3].Category, whitelist[3].Id).
		AddRow(user.Id, whitelist[4].Name, whitelist[4].Category, whitelist[4].Id).
		AddRow(user.Id, whitelist[5].Name, whitelist[5].Category, whitelist[5].Id)

	query2 := fmt.Sprintf("SELECT id, name, category, itemId FROM whitelist WHERE whitelist.id \\= %q", user.Id)
	mock.
		ExpectQuery(query2).
		WillReturnRows(rows2)
	
	b, err := models.GetBudget(app, user.Id)
	test.ModelMethod(t, err, "select")
	test.MockExpectations(t, mock)

	if (b.Categories[0].Id != budget.Categories[0].Id && 
		b.Categories[0].WhiteList[0].Id != budget.Categories[0].WhiteList[0].Id) {
		t.Error("The function successfully executed but there was an error getting the correct budget")
	}
}