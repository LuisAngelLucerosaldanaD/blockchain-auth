package users_temp

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de UsersTemp
type UserTemp struct {
	ID           string    `json:"id" db:"id" valid:"required,uuid"`
	Nickname     string    `json:"nickname" db:"nickname" valid:"required"`
	Email        string    `json:"email" db:"email" valid:"required"`
	Password     string    `json:"password" db:"password" valid:"required"`
	Name         string    `json:"name" db:"name" valid:"required"`
	Lastname     string    `json:"lastname" db:"lastname" valid:"required"`
	IdType       int       `json:"id_type" db:"id_type" valid:"required"`
	IdNumber     string    `json:"id_number" db:"id_number" valid:"required"`
	Cellphone    string    `json:"cellphone" db:"cellphone" valid:"required"`
	BirthDate    time.Time `json:"birth_date" db:"birth_date" valid:"required"`
	VerifiedCode string    `json:"verified_code" db:"verified_code" valid:"required"`
	IsDeleted    bool      `json:"is_deleted" db:"is_deleted"`
	IdUser       string    `json:"id_user" db:"id_user"`
	DeletedAt    time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func NewUsersTemp(id string, nickname string, email string, password string, name string, lastname string, idType int, idNumber string, cellphone string, birthDate time.Time, VerifiedCode string) *UserTemp {
	return &UserTemp{
		ID:           id,
		Nickname:     nickname,
		Email:        email,
		Password:     password,
		Name:         name,
		Lastname:     lastname,
		IdType:       idType,
		IdNumber:     idNumber,
		Cellphone:    cellphone,
		BirthDate:    birthDate,
		VerifiedCode: VerifiedCode,
		IsDeleted:    false,
		IdUser:       id,
	}
}

func (m *UserTemp) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
