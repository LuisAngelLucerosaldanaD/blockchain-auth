package wallets_temp

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de WalletTemp
type WalletTemp struct {
	ID        string    `json:"id" db:"id" valid:"required,uuid"`
	Mnemonic  string    `json:"mnemonic" db:"mnemonic" valid:"required"`
	IdUser    string    `json:"id_user" db:"id_user" valid:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewWalletTemp(id string, mnemonic string, idUser string) *WalletTemp {
	return &WalletTemp{
		ID:       id,
		Mnemonic: mnemonic,
		IdUser:   idUser,
	}
}

func (m *WalletTemp) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
