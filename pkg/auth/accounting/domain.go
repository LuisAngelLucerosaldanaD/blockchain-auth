package accounting

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de Accounting
type Accounting struct {
	ID        string    `json:"id" db:"id" valid:"uuid"`
	IdWallet  string    `json:"id_wallet" db:"id_wallet" valid:"required"`
	Amount    float64   `json:"amount" db:"amount"`
	IdUser    string    `json:"id_user" db:"id_user" valid:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewAccounting(id string, idWallet string, amount float64, idUser string) *Accounting {
	return &Accounting{
		ID:       id,
		IdWallet: idWallet,
		Amount:   amount,
		IdUser:   idUser,
	}
}

func NewAccountingSetAmount(idWallet string, amount float64, idUser string) *Accounting {
	return &Accounting{
		IdWallet: idWallet,
		Amount:   amount,
		IdUser:   idUser,
	}
}

func (m *Accounting) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
