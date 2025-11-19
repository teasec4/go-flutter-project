package account

type Account interface {
	GetBalance() float64
}

type impl struct {
	AccountId string
	Name      string
	Balance   float64
}

var store = map[string]Account{}

func NewAccount(accountId string, name string, balance float64) Account {
	a := &impl{
		AccountId: accountId,
		Name:      name,
		Balance:   balance,
	}
	store[accountId] = a
	return a
}

func GetByID(accountId string) (Account, bool) {
	a, ok := store[accountId]
	return a, ok
}

func (a *impl) GetBalance() float64 {
	return a.Balance
}
