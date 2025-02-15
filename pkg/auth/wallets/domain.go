package wallets

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Wallet estructura de la wallet
type Wallet struct {
	ID        string    `json:"id" db:"id" valid:"required,uuid"`
	Mnemonic  string    `json:"mnemonic" db:"mnemonic" valid:"required"`
	RsaPublic string    `json:"rsa_public" db:"rsa_public" valid:"required"`
	IpDevice  string    `json:"ip_device" db:"ip_device" valid:"required"`
	StatusId  int       `json:"status_id" db:"status_id" valid:"required"`
	Dni       string    `json:"dni" db:"dni" valid:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewWallet(id, mnemonic, rsaPublic, ipDevice, identityNumber string, statusId int) *Wallet {
	return &Wallet{
		ID:        id,
		Mnemonic:  mnemonic,
		RsaPublic: rsaPublic,
		IpDevice:  ipDevice,
		StatusId:  statusId,
		Dni:       identityNumber,
	}
}

func UpdateWallet(id, mnemonic, ipDevice, identityNumber string, statusId int) *Wallet {
	return &Wallet{
		ID:       id,
		Mnemonic: mnemonic,
		IpDevice: ipDevice,
		StatusId: statusId,
		Dni:      identityNumber,
	}
}

func (m *Wallet) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
