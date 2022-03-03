package login

import (
	"blion-auth/internal/ciphers"
	"blion-auth/internal/logger"
	"blion-auth/internal/msg"
	"blion-auth/pkg/auth"
	b64 "encoding/base64"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type handlerLogin struct {
	DB   *sqlx.DB
	TxID string
}

// Login godoc
// @Summary blockchain
// @Description GetBlockToMine
// @Accept  json
// @Produce  json
// @Success 200 {object} responseBlockToMine
// @Success 202 {object} dataBlockToMine
// @Router /api/v1/block-to-mine [get]
// @Authorization Bearer token
func (h *handlerLogin) Login(c *fiber.Ctx) error {

	res := responseLogin{Error: true}
	m := requestLogin{}
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf("couldn't bind model login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	token, cod, err := srvAuth.SrvLogin.Login(m.Nickname, m.Email, m.Password, c.IP())
	if err != nil {
		logger.Warning.Printf("couldn't login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(cod, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	authRes := Token{AccessToken: token, RefreshToken: token}
	res.Data = authRes
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

func (h *handlerLogin) SecretKey(c *fiber.Ctx) error {

	res := responseKey{Error: true}
	if c.Params("secret") != "027dfc14-d847-4627-9f7f-ecb5d6ef06fa" {
		res.Data = ""
		res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
		res.Error = false
		return c.Status(http.StatusOK).JSON(res)
	}

	res.Data = b64.StdEncoding.EncodeToString([]byte(ciphers.GetSecret()))
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
