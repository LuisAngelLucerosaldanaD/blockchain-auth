package users_temp

import (
	"github.com/jmoiron/sqlx"

	"blion-auth/internal/logger"
	"blion-auth/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesUsersTempRepository interface {
	create(m *UserTemp) error
	update(m *UserTemp) error
	delete(id string) error
	getByID(id string) (*UserTemp, error)
	getAll() ([]*UserTemp, error)
	getByEmail(email string) (*UserTemp, error)
	getByNickname(nickname string) (*UserTemp, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesUsersTempRepository {
	var s ServicesUsersTempRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newUsersTempPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
