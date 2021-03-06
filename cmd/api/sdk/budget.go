package sdk

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/elopez00/scale-backend/cmd/api/models"
	"github.com/elopez00/scale-backend/pkg/application"
	"github.com/julienschmidt/httprouter"
)

// UpdateBudget handler will be in charge of all forms of edits whether it be additions, deletions,
// or changes by using a request body in the format of an UpdateRequest. This function will
// be used to create and change the user's budget in the application. If there is an error,
// with a query or database connection it will be logged and returned as a JSON response.
func UpdateBudget(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer CloseBody(r)

		// extract the budget information from the request body
		var budget models.Budget
		err := json.NewDecoder(r.Body).Decode(&budget)
		if err != nil {
			log.Println("Failed to decode budget from body")
			return
		}

		// gets the user id extracted from authentication cookie for later
		// use in the creation of the row containing the permanent token
		userId := fmt.Sprintf("%v", r.Context().Value(models.Key("user")))

		// updates any items in the budget.Request
		if err := budget.Update(app, userId); err != nil {
			msg := "Failed to store budget information in database"
			models.CreateError(w, http.StatusBadGateway, msg, err)
			return
		}

		msg := "Successfully created budget"
		models.CreateResponse(w, msg, nil)
	}
}

// GetBudget handler gets the budget from the databse in the form of a JSONified budget
// object returned in the responses result property. If there is an error with the
// database connection or query it will be logged nad returned as a JSON response.
// TODO make and test this function pls thank you!
func GetBudget(app *application.App) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	}
}
