package wallets

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de Wallet
type Wallet struct {
	ID               string    `json:"id" db:"id" valid:"required,uuid"`
	Mnemonic         string    `json:"mnemonic" db:"mnemonic" valid:"required"`
	RsaPublic        string    `json:"rsa_public" db:"rsa_public" valid:"required"`
	RsaPrivate       string    `json:"rsa_private" db:"rsa_private" valid:"required"`
	RsaPublicDevice  string    `json:"rsa_public_device" db:"rsa_public_device" valid:"required"`
	RsaPrivateDevice string    `json:"rsa_private_device" db:"rsa_private_device" valid:"required"`
	IpDevice         string    `json:"ip_device" db:"ip_device" valid:"required"`
	StatusId         int       `json:"status_id" db:"status_id" valid:"required"`
	IdentityNumber   string    `json:"identity_number" db:"identity_number" valid:"required"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

func NewWallet(id, mnemonic, rsaPublic, rsaPrivate, rsaPublicDevice, rsaPrivateDevice, ipDevice, identityNumber string, statusId int) *Wallet {
	return &Wallet{
		ID:               id,
		Mnemonic:         mnemonic,
		RsaPublic:        rsaPublic,
		RsaPrivate:       rsaPrivate,
		RsaPublicDevice:  rsaPublicDevice,
		RsaPrivateDevice: rsaPrivateDevice,
		IpDevice:         ipDevice,
		StatusId:         statusId,
		IdentityNumber:   identityNumber,
	}
}

func (m *Wallet) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
