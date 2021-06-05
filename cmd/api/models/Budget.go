package models

import (
	application "github.com/elopez00/scale-backend/pkg/app"
)

// item for white listed companies of a specific category
type WhiteListItem struct {
	Category string `json:"category"` // name of category
	Name     string `json:"name"`     // name of company being listed under this category
	Id		 string `json:"id"`       // item id
}

// category within budget
type Category struct {
	Name   string  `json:"name"`   // name of category
	Budget float64 `json:"budget"` // amount of money budgeted towards this category
	Id	   string  `json:"id"`	   // category id
}

// Object containing both category updates and whitelist updates. Neither one or the
// other are required.
type UpdateObject struct {
	Categories []Category  `json:"categories,omitempty"`
	WhiteList  []WhiteListItem `json:"whitelist,omitempty"`
}

// a struct that describes what will be updated in a budget
type UpdateRequest struct {
	Update UpdateObject `json:"change,omitempty"`
	Remove UpdateObject `json:"remove,omitempty"`
}

// A budget is a combination of categories
type Budget struct {
	Categories []Category      `json:"categories"`
	WhiteList  []WhiteListItem `json:"whitelist,omitempty"`
	Request    UpdateRequest   `json:"request,omitempty"`
}

// Single function that handles all updates to current budget whether it be adding, removing,
// or changing. This function will only perform at most 4 queries at a time. If there is a
// failure inserting, deleting, or updating any of the rows it will be returned as an error.
func (b *Budget) Update(app *application.App, userId string) error {
	// add any categories that need to be added
	if err := UpdateCategories(app, userId, b.Request.Update.Categories); err != nil {
		return err
	}

	// add any whitelist elements that need to be added
	if err := UpdateWhiteList(app, userId, b.Request.Update.WhiteList); err != nil {
		return err
	}

	return nil
}

// Gets all the white list items and inserts them to the database. If the function fails
// due to the databse connection or query execution, an error will be returned that reflects
// this
func UpdateWhiteList(app *application.App, userId string, whitelist []WhiteListItem) error {
	// there might not be items that needs to be whitelisted, if this is the case return nil
	if len(whitelist) == 0 {
		return nil
	}

	query := "INSERT INTO whitelist(id, category, name, itemId) VALUES "
	queryEnd := 
		" AS updated ON DUPLICATE KEY UPDATE id=updated.id, category=updated.category,"+
		" name=updated.name, itemId=updated.itemId;"

	vals := []interface{}{}

	for _, item := range whitelist {
		query += " (?,?,?,?),"
		vals = append(vals, userId, item.Category, item.Name, item.Id)
	}

	// prepare statement
	query = query[0 : len(query)-1] + queryEnd
	stmt, err := app.DB.Client.Prepare(query)
	if err != nil {
		return err
	}

	// execute query
	if _, err := stmt.Exec(vals...); err != nil {
		return err
	}

	return nil
}

// Gets all the category items and inserts them to the database. If the function fails due
// to the databse connection or query execution, an error will be returned that reflects this
func UpdateCategories(app *application.App, userId string, categories []Category) error {
	if len(categories) == 0 {
		return nil
	}

	query := "INSERT INTO categories(id, name, budget, categoryId) VALUES "
	queryEnd := 
		" AS updated ON DUPLICATE KEY UPDATE id=updated.id, name=updated.name,"+
		" budget=updated.budget, categoryId=updated.categoryId;"

	vals := []interface{}{}

	for _, category := range categories {
		query += " (?,?,?,?),"
		vals = append(vals, userId, category.Name, category.Budget, category.Id)
	}

	// prepare statement
	query = query[0 : len(query)-1] + queryEnd// trim last comma
	stmt, err := app.DB.Client.Prepare(query)
	if err != nil {
		return err
	}

	// execute query
	if _, err := stmt.Exec(vals...); err != nil {
		return err
	}

	return nil
}
