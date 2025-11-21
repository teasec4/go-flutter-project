package tool

import "time"

type mockDB struct {}

var mockLogingDetails = map[string]LoginDetails{
	"1234567890": {
		AuthToken: "1234567890",
		AccountId: "test",
	},
	"2345678901": {
		AuthToken: "2345678901",
		AccountId: "test2",
	},
}

var mockAccountDetails = map[string]AccountDetails{
	"1234567890": {	
		Balance: 1000,
		AccountId: "test",
	},
	"2345678901": {
		Balance: 2000,
		AccountId: "test2",
	},
}

func (d *mockDB) GetUserBalance(accountId string) *AccountDetails {
	time.Sleep(100 * time.Millisecond)
	clientData, ok := mockAccountDetails[accountId]
	if !ok {
		return nil
	}
	return &clientData
}

func (d *mockDB) GetUserLoginDetails(accountId string) *LoginDetails {
	time.Sleep(100 * time.Millisecond)
	clientData, ok := mockLogingDetails[accountId]
	if !ok {
		return nil
	}
	return &clientData
}

func (d *mockDB) SetupDatabase() error {
	return nil
}
