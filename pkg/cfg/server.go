package cfg

import (
	"blion-auth/internal/models"
	"blion-auth/pkg/cfg/dictionaries"
	"blion-auth/pkg/cfg/messages"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	SrvDictionaries dictionaries.PortsServerDictionaries
	SrvMessage      messages.PortsServerMessages
}

func NewServerCfg(db *sqlx.DB, user *models.User, txID string) *Server {

	repoDictionaries := dictionaries.FactoryStorage(db, user, txID)
	srvDictionaries := dictionaries.NewDictionariesService(repoDictionaries, user, txID)

	repoMessage := messages.FactoryStorage(db, user, txID)
	srvMessage := messages.NewMessagesService(repoMessage, user, txID)

	return &Server{
		SrvDictionaries: srvDictionaries,
		SrvMessage:      srvMessage,
	}
}
