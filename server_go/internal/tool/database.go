package tool

import "log"

type LoginDetails struct{
	AuthToken string
	AccountId string
}

type AccountDetails struct{
	Balance int64
	AccountId string
}

type DatabaseInterface interface{
	GetUserLoginDetails(accountId string) *LoginDetails
	GetUserBalance(accountId string) *AccountDetails
	SetupDatabase()error
}

func NewDatabase() (DatabaseInterface, error) {
	var database DatabaseInterface = &mockDB{}

	err := database.SetupDatabase()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return database, nil
}
