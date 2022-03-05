package users

import (
	"blion-auth/internal/email"
	"blion-auth/internal/env"
	"blion-auth/internal/grpc/users_proto"
	"blion-auth/internal/logger"
	"blion-auth/internal/models"
	"blion-auth/internal/msg"
	"blion-auth/internal/password"
	"blion-auth/pkg/auth"
	"blion-auth/pkg/auth/login"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"strconv"
	"time"
)

type HandlerUsers struct {
	DB   *sqlx.DB
	TxID string
}

func (h HandlerUsers) CreateUser(ctx context.Context, request *users_proto.UserRequest) (*users_proto.UserResponse, error) {
	res := &users_proto.UserResponse{Error: true}

	var parameters = make(map[string]string, 0)
	e := env.NewConfiguration()

	//TODO implements encrypt-decrypt password
	if request.Password != request.ConfirmPassword {
		logger.Error.Printf("this password is not equal to confirm_password")
		res.Code, res.Type, res.Msg = msg.GetByCode(10005, h.DB, h.TxID)
		return res, fmt.Errorf("this password is not equal to confirm_password")
	}

	srvUser := auth.NewServerAuth(h.DB, nil, h.TxID)
	id := uuid.New().String()

	min := 1000
	max := 9999
	rand.Seed(time.Now().UnixNano())
	emailCode := strconv.Itoa(rand.Intn(max-min+1) + min)
	verifiedCode := password.Encrypt(emailCode)

	request.Password = password.Encrypt(request.Password)

	// TODO valid password and time
	usr, code, err := srvUser.SrvUserTemp.CreateUsersTemp(id, request.Nickname, request.Email, request.Password, request.Name, request.Lastname,
		int(request.IdType), request.IdNumber, request.Cellphone, time.Now(), verifiedCode)
	if err != nil {
		logger.Error.Printf("couldn't create User: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	usr.Password = ""
	usr.VerifiedCode = ""

	parameters["@access-code"] = emailCode
	tos := []string{request.Email}

	err = email.Send(tos, parameters, e.Template.EmailCode, "Access code to verify user")
	if err != nil {
		logger.Error.Println("error when execute send email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(10002, h.DB, h.TxID)
		return res, err
	}

	usrTemp := models.User{
		ID:        usr.ID,
		Nickname:  usr.Nickname,
		Email:     usr.Email,
		Name:      usr.Name,
		Lastname:  usr.Lastname,
		IdType:    usr.IdType,
		IdNumber:  usr.IdNumber,
		Cellphone: usr.Cellphone,
		BirthDate: usr.BirthDate,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
	}

	token, code, err := login.GenerateJWT(&usrTemp)

	if err != nil {
		logger.Error.Println("error, don't create token: %V", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	parameters["@url-token"] = e.Portal.Url + e.Portal.ActivateAccount + token

	err = email.Send(tos, parameters, e.Template.EmailToken, "Verify account")
	if err != nil {
		logger.Error.Println("error when execute send email: %V", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(10002, h.DB, h.TxID)
		return res, err
	}

	res.Data = "usr"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return  res, nil

}

func (h HandlerUsers) ActivateUser(ctx context.Context, request *users_proto.ActivateUserRequest) (*users_proto.ActivateUserResponse, error) {
	panic("implement me")
}

func (h HandlerUsers) ValidateEmail(ctx context.Context, request *users_proto.ValidateEmailRequest) (*users_proto.ValidateEmailResponse, error) {
	panic("implement me")
}

func (h HandlerUsers) ValidateNickname(ctx context.Context, request *users_proto.ValidateNicknameRequest) (*users_proto.ValidateNicknameResponse, error) {
	panic("implement me")
}

func (h HandlerUsers) GetUserById(ctx context.Context, request *users_proto.GetUserByIDRequest) (*users_proto.UserResponse, error) {
	panic("implement me")
}

func (h HandlerUsers) ValidateIdentity(ctx context.Context, request *users_proto.ValidateIdentityRequest) (*users_proto.ValidateIdentityResponse, error) {
	panic("implement me")
}

func (h HandlerUsers) ValidateCertifier(ctx context.Context, request *users_proto.ValidateCertifierRequest) (*users_proto.ValidateCertifierResponse, error) {
	panic("implement me")
}


