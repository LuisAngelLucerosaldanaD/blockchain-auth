package dictionaries

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de Dictionaries
type Dictionaries struct {
	ID          int       `json:"id" db:"id" valid:"-"`
	Name        string    `json:"name" db:"name" valid:"required"`
	Value       string    `json:"value" db:"value" valid:"required"`
	Description string    `json:"description" db:"description" valid:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewDictionaries(id int, name string, value string, description string) *Dictionaries {
	return &Dictionaries{
		ID:          id,
		Name:        name,
		Value:       value,
		Description: description,
	}
}

func NewCreateDictionaries(name string, value string, description string) *Dictionaries {
	return &Dictionaries{
		Name:        name,
		Value:       value,
		Description: description,
	}
}

func (m *Dictionaries) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
