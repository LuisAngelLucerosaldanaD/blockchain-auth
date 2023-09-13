package auth

import (
	"blion-auth/internal/models"
	"blion-auth/pkg/auth/accounting"
	"blion-auth/pkg/auth/frozen_money"
	"blion-auth/pkg/auth/wallets"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	SrvWallet      wallets.PortsServerWallet
	SrvAccounting  accounting.PortsServerAccounting
	SrvFrozenMoney frozen_money.PortsServerFrozenMoney
}

func NewServerAuth(db *sqlx.DB, user *models.User, txID string) *Server {
	repoWallet := wallets.FactoryStorage(db, user, txID)
	srvWallet := wallets.NewWalletService(repoWallet, user, txID)
	repoAccounting := accounting.FactoryStorage(db, user, txID)
	srvAccounting := accounting.NewAccountingService(repoAccounting, user, txID)

	repoFrozen := frozen_money.FactoryStorage(db, user, txID)
	srvFrozen := frozen_money.NewFrozenMoneyService(repoFrozen, user, txID)

	return &Server{
		SrvWallet:      srvWallet,
		SrvAccounting:  srvAccounting,
		SrvFrozenMoney: srvFrozen,
	}
}
