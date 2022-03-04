package msg

import (
	"blion-auth/internal/env"
	"blion-auth/internal/logger"
	"blion-auth/pkg/cfg"
	"strconv"

	"github.com/jmoiron/sqlx"
)

func GetByCode(code int, db *sqlx.DB, txID string) (int32, int32, string) {
	var codRes int32
	msg := ""
	c := env.NewConfiguration()
	srvCFG := cfg.NewServerCfg(db, nil, txID)
	m, codErr, err := srvCFG.SrvMessage.GetMessagesByID(code)
	if err != nil {
		return codRes, 0, strconv.Itoa(codErr)
	}

	switch c.App.Language {
	case "sp":
		msg = m.Spa
	case "en":
		msg = m.Eng
	default:
		logger.Error.Println("el sistema no tiene implementado el idioma: ", c.App.Language)
	}
	codRes = int32(m.ID)
	return codRes, int32(m.TypeMessage), msg
}
