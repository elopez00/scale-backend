package models

type User struct {
	Id			string	`json:"id,omitempty"`
	FirstName	string	`json:"firstname,omitempty"`
	LastName	string	`json:"lastname,omitempty"`
	Email		string	`json:"email,omitempty"`
	Password	string	`json:"password,omitempty"`
}

type Response struct {
	Status		int		`json:"status"`
	Message		string	`json:"message"`
	Type 		string 	`json:"type"`
}

type Env struct {
	DbUser		string
	DbPass		string
	DbData		string
	DbAccess	string
	PlaidCtry	string
	PlaidSecret	string
	PlaidId		string
	PlaidRedir	string
	Live		string
	Key 		string
}
