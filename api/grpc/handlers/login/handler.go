package login

import (
	"blion-auth/internal/ciphers"
	"blion-auth/internal/grpc/auth_proto"
	"blion-auth/internal/logger"
	"blion-auth/internal/msg"
	"blion-auth/pkg/auth"
	"context"
	"encoding/base64"
	"github.com/jmoiron/sqlx"
)

type HandlerLogin struct {
	DB   *sqlx.DB
	TxID string
}

func (h *HandlerLogin) Login(ctx context.Context, request *auth_proto.LoginRequest) (*auth_proto.LoginResponse, error) {
	res := auth_proto.LoginResponse{Error: true}

	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	var nickName, email string
	if request.Nickname != nil {
		nickName = *request.Nickname
	}
	if request.Email != nil {
		email = *request.Email
	}
	token, cod, err := srvAuth.SrvLogin.Login(nickName, email, request.Password, "")
	if err != nil {
		logger.Warning.Printf("couldn't login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(cod, h.DB, h.TxID)
		return &res, err
	}
	tkn := &auth_proto.DataMsg{
		AccessToken:  token,
		RefreshToken: token,
	}
	res.Data = tkn
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return &res, err
}

func (h *HandlerLogin) SecretKey(ctx context.Context, request *auth_proto.SecretRequest) (*auth_proto.SecretResponse, error)  {
	res := auth_proto.SecretResponse{Error: true}
	if request.Secret != "027dfc14-d847-4627-9f7f-ecb5d6ef06fa" {
		res.Data = ""
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return &res, nil
	}
	res.Data = base64.StdEncoding.EncodeToString([]byte(ciphers.GetSecret()))
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return &res, nil
}
