package auth

import (
	"blion-auth/internal/models"
	"blion-auth/pkg/auth/accounting"
	"blion-auth/pkg/auth/login"
	"blion-auth/pkg/auth/users"
	"blion-auth/pkg/auth/users_temp"
	"blion-auth/pkg/auth/users_wallet"
	"blion-auth/pkg/auth/wallets"
	"blion-auth/pkg/auth/wallets_temp"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	SrvLogin       login.PortServiceLogin
	SrvUserTemp    users_temp.PortsServerUsersTemp
	SrvUser        users.PortsServerUsers
	SrvWalletTemp  wallets_temp.PortsServerWalletTemp
	SrvWallet      wallets.PortsServerWallet
	SrvAccounting  accounting.PortsServerAccounting
	SrvUsersWallet users_wallet.PortsServerUsersWallet
}



func NewServerAuth(db *sqlx.DB, user *models.User, txID string) *Server {
	srvLogin := login.NewLoginService(db, txID)
	repoUserTemp := users_temp.FactoryStorage(db, user, txID)
	srvUserTemp := users_temp.NewUsersTempService(repoUserTemp, user, txID)
	repoUser := users.FactoryStorage(db, user, txID)
	srvUser := users.NewUsersService(repoUser, user, txID)
	repoWalletTemp := wallets_temp.FactoryStorage(db, user, txID)
	srvWalletTemp := wallets_temp.NewWalletTempService(repoWalletTemp, user, txID)
	repoWallet := wallets.FactoryStorage(db, user, txID)
	srvWallet := wallets.NewWalletService(repoWallet, user, txID)
	repoAccounting := accounting.FactoryStorage(db, user, txID)
	srvAccounting := accounting.NewAccountingService(repoAccounting, user, txID)

	repoUsersWallet := users_wallet.FactoryStorage(db, user, txID)
	srvUsersWallet := users_wallet.NewUsersWalletService(repoUsersWallet, user, txID)

	return &Server{
		SrvLogin:       srvLogin,
		SrvUserTemp:    srvUserTemp,
		SrvUser:        srvUser,
		SrvWalletTemp:  srvWalletTemp,
		SrvWallet:      srvWallet,
		SrvAccounting:  srvAccounting,
		SrvUsersWallet: srvUsersWallet,
	}
}
