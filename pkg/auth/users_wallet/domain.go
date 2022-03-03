package users_wallet

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de Wallet
type UsersWallet struct {
	ID        string    `json:"id" db:"id" valid:"required"`
	IdUser    string    `json:"id_user" db:"id_user" valid:"required"`
	IdWallet  string    `json:"id_wallet" db:"id_wallet" valid:"required"`
	IsDelete  bool      `json:"is_delete" db:"is_delete" valid:"-"`
	DeletedAt time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewWallet(id, idUser, idWallet string, isDelete bool) *UsersWallet {
	return &UsersWallet{
		ID:       id,
		IdUser:   idUser,
		IdWallet: idWallet,
		IsDelete: isDelete,
	}
}

func (m *UsersWallet) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
