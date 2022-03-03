package wallets_temp

import (
	"github.com/jmoiron/sqlx"

	"blion-auth/internal/logger"
	"blion-auth/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesWalletTempRepository interface {
	create(m *WalletTemp) error
	update(m *WalletTemp) error
	delete(id string) error
	getByID(id string) (*WalletTemp, error)
	getAll() ([]*WalletTemp, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesWalletTempRepository {
	var s ServicesWalletTempRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newWalletTempPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
