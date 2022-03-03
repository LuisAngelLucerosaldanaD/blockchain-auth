package accounting

import (
	"github.com/jmoiron/sqlx"

	"blion-auth/internal/logger"
	"blion-auth/internal/models"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

type ServicesAccountingRepository interface {
	create(m *Accounting) error
	update(m *Accounting) error
	delete(id string) error
	getByID(id string) (*Accounting, error)
	getAll() ([]*Accounting, error)
	setAmount(m *Accounting) error
	getByWalletID(walletID string) (*Accounting, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesAccountingRepository {
	var s ServicesAccountingRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newAccountingPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
