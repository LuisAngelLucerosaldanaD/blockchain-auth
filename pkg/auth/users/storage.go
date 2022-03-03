package users

import (
	"github.com/jmoiron/sqlx"

	"blion-auth/internal/logger"
	"blion-auth/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesUsersRepository interface {
	create(m *User) error
	update(m *User) error
	delete(id string) error
	getByID(id string) (*User, error)
	getAll() ([]*User, error)
	getByEmail(email string) (*User, error)
	getByNickname(nickname string) (*User, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesUsersRepository {
	var s ServicesUsersRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newUsersPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
