package users_wallet

import (
	"github.com/jmoiron/sqlx"

	"blion-auth/internal/logger"
	"blion-auth/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesUsersWalletRepository interface {
	create(m *UsersWallet) error
	update(m *UsersWallet) error
	delete(id string) error
	getByID(id string) (*UsersWallet, error)
	getAll() ([]*UsersWallet, error)
	getByUserIDAndIdentityNumber(userID, identityNumber string) (*UsersWallet, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesUsersWalletRepository {
	var s ServicesUsersWalletRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newWalletPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
