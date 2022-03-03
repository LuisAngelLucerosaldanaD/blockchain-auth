package login

import (
	"blion-auth/internal/grpc/auth_proto"
	"blion-auth/internal/logger"
	"blion-auth/internal/models"
	"blion-auth/internal/pwd"
	"blion-auth/pkg/auth/users"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PortServiceLogin interface {
	Login(context.Context, *auth_proto.LoginRequest) (*auth_proto.Response, error)
}

type service struct {
	DB   *sqlx.DB
	TxID string
}



func NewLoginService(db *sqlx.DB, txID string) PortServiceLogin {
	return &service{DB: db, TxID: txID}
}
func (s *service) Login(ctx context.Context, request *auth_proto.LoginRequest) (*auth_proto.Response, error) {
	panic("implement me")
}
func (s *service) Login2(nickname, email, password, realIp string) (string, int, error) {
	var token string
	m := NewLogin(nickname, email, password, realIp)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.TxID, " - don't meet validations:", err)
		return "", 15, err
	}
	usr, cod, err := s.getUserByEmailOrNickname(nickname, email)
	if err != nil {
		logger.Error.Println(s.TxID, " - couldn't get wallets by id", err)
		return token, cod, err
	}
	if !pwd.Compare(usr.ID, usr.Password, m.Password) {
		return token, 10, fmt.Errorf("wallets o password incorrecto")
	}
	token, cod, err = GenerateJWT(s.resetFieldsUser(usr))
	if err != nil {
		logger.Error.Printf("couldn't get token:%s - %v", s.TxID, err)
		return "", cod, err
	}
	return token, 29, nil
}

func (s *service) getUserByEmailOrNickname(nickname, email string) (*models.User, int, error) {
	var usr models.User
	repoUsers := users.FactoryStorage(s.DB, nil, s.TxID)
	srvUsers := users.NewUsersService(repoUsers, nil, s.TxID)

	if email != "" {
		user, _, err := srvUsers.GetUsersByEmail(email)
		if err != nil {
			logger.Error.Println("couldn't get user by id", err)
			return nil, 10, err
		}
		if user == nil {
			return nil, 10, fmt.Errorf("user or password not found")
		}
		usr = models.User(*user)
	}
	if nickname != "" {
		user, _, err := srvUsers.GetUsersByNickname(nickname)
		if err != nil {
			logger.Error.Println("user or password not found", err)
			return nil, 10, err
		}
		if user == nil {
			return nil, 10, fmt.Errorf("user or password not found")
		}
		usr = models.User(*user)
	}
	if &usr == nil {
		return nil, 10, fmt.Errorf("couldn't get user by id")
	}

	return &usr, 29, nil
}

func (s *service) resetFieldsUser(user *models.User) *models.User {
	//TODO: reset fields date and verify code.
	user.Password = ""
	user.RsaPrivate = ""
	user.RsaPublic = ""
	user.FailedAttempts = 0
	user.IdType = 0
	user.VerifiedCode = ""
	user.IsDeleted = false
	user.IdUser = ""
	user.FullPathPhoto = ""
	return user
}
