package wallets

import (
	"github.com/jmoiron/sqlx"

	"blion-auth/internal/logger"
	"blion-auth/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesWalletRepository interface {
	create(m *Wallet) error
	update(m *Wallet) error
	delete(id string) error
	getByID(id string) (*Wallet, error)
	getAll() ([]*Wallet, error)
	getWalletByUserId(userID string) ([]*Wallet, error)
	getWalletByIdentityNumber(identityNumber string) (*Wallet, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesWalletRepository {
	var s ServicesWalletRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newWalletPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
