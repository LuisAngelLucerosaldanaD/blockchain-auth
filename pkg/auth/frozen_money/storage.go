package frozen_money

import (
	"github.com/jmoiron/sqlx"

	"blion-auth/internal/logger"
	"blion-auth/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesFrozenMoneyRepository interface {
	create(m *FrozenMoney) error
	update(m *FrozenMoney) error
	delete(id string) error
	getByID(id string) (*FrozenMoney, error)
	getAll() ([]*FrozenMoney, error)
	getByWalletIDAndLotteryId(walletID, lotteryId string) (*FrozenMoney, error)
	getByWalletID(walletID string) (*FrozenMoney, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesFrozenMoneyRepository {
	var s ServicesFrozenMoneyRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newFrozenMoneyPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
