package dictionaries

import (
	"github.com/jmoiron/sqlx"

	"blion-auth/internal/logger"
	"blion-auth/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesDictionariesRepository interface {
	create(m *Dictionaries) error
	update(m *Dictionaries) error
	delete(id int) error
	getByID(id int) (*Dictionaries, error)
	getAll() ([]*Dictionaries, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesDictionariesRepository {
	var s ServicesDictionariesRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newDictionariesPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
