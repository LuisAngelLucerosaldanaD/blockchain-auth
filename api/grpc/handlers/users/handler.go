package users

import (
	"blion-auth/internal/email"
	"blion-auth/internal/env"
	"blion-auth/internal/grpc/users_proto"
	"blion-auth/internal/helpers"
	"blion-auth/internal/logger"
	"blion-auth/internal/mnemonic"
	"blion-auth/internal/models"
	"blion-auth/internal/msg"
	"blion-auth/internal/password"
	"blion-auth/internal/pwd"
	"blion-auth/internal/rsa_generate"
	"blion-auth/pkg/auth"
	"blion-auth/pkg/auth/login"
	"blion-auth/pkg/auth/users"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"strconv"
	"strings"
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

	if request.Password != request.ConfirmPassword {
		logger.Error.Printf("this password is not equal to confirm_password")
		res.Code, res.Type, res.Msg = msg.GetByCode(10005, h.DB, h.TxID)
		return res, fmt.Errorf("this password is not equal to confirm_password")
	}

	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	usrI, code, err := srvAuth.SrvUser.GetUserByIdentityNumber(request.IdNumber)
	if err != nil {
		logger.Error.Printf("couldn't get user by identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return res, err
	}

	if usrI != nil {
		logger.Error.Printf("Ya exisite un usuario con el numero de identificación ingresado: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return res, fmt.Errorf("ya exisite un usuario con el numero de identificación ingresado")
	}

	id := uuid.New().String()

	min := 1000
	max := 9999
	rand.Seed(time.Now().UnixNano())
	emailCode := strconv.Itoa(rand.Intn(max-min+1) + min)
	verifiedCode := password.Encrypt(emailCode)

	request.Password = password.Encrypt(request.Password)
	layout := "2006-01-02T15:04:05.000Z"
	var birthDate time.Time
	birthDate, err = time.Parse(layout, request.BirthDate)
	if err != nil {
		birthDate = time.Now()
	}
	usr, code, err := srvAuth.SrvUserTemp.CreateUsersTemp(id, request.Nickname, request.Email, request.Password,
		request.Name, request.Lastname, int(request.IdType), request.IdNumber, request.Cellphone,
		birthDate, verifiedCode)
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
	userResponse := &users_proto.User{
		ID:        usr.ID,
		Nickname:  usr.Nickname,
		Email:     usr.Email,
		Name:      usr.Name,
		Lastname:  usr.Lastname,
		IdType:    int32(usr.IdType),
		IdNumber:  usr.IdNumber,
		Cellphone: usr.Cellphone,
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

	request.Id = id
	request.Password = ""
	request.ConfirmPassword = ""
	res.Data = userResponse
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false

	return res, nil

}

func (h HandlerUsers) ActivateUser(ctx context.Context, request *users_proto.ActivateUserRequest) (*users_proto.ValidateResponse, error) {
	res := &users_proto.ValidateResponse{Error: true}
	u, err := helpers.GetUserContext(ctx)
	if err != nil {
		logger.Error.Printf("couldn't get token user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return res, err
	}

	rsaPrivate, rsaPublic, err := rsa_generate.Execute()
	if err != nil {
		logger.Error.Printf("couldn't generate rsa user in ActivateUser: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return res, err
	}

	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	usrTemp, code, err := srvAuth.SrvUserTemp.GetUsersTempByID(u.ID)
	if err != nil {
		logger.Error.Printf("couldn't get user by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	if !password.Compare(usrTemp.ID, usrTemp.VerifiedCode, request.Code) {
		logger.Error.Printf("the verification code is not correct: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(10005, h.DB, h.TxID)
		return res, err
	}

	usr, code, err := srvAuth.SrvUser.CreateUsers(usrTemp.ID, usrTemp.Nickname, usrTemp.Email, usrTemp.Password,
		usrTemp.Name, usrTemp.Lastname, usrTemp.IdType, usrTemp.IdNumber, usrTemp.Cellphone, usrTemp.BirthDate,
		"", time.Now(), "", rsaPrivate, rsaPublic, 21)
	if err != nil {
		logger.Error.Printf("don't create user: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	code, err = srvAuth.SrvUserTemp.DeleteUsersTemp(usr.ID)
	if err != nil {
		logger.Error.Printf("don't delete user temp by email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	res.Data = "Active Account Successful"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h HandlerUsers) ValidateEmail(ctx context.Context, request *users_proto.ValidateEmailRequest) (*users_proto.ValidateResponse, error) {
	res := &users_proto.ValidateResponse{Error: true}

	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	usrTemp, code, err := srvAuth.SrvUserTemp.GetUserByEmail(request.Email)
	if err != nil {
		logger.Error.Printf("couldn't get user by email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	usr, code, err := srvAuth.SrvUser.GetUsersByEmail(request.Email)
	if err != nil {
		logger.Error.Printf("couldn't get user by email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	if usrTemp != nil || usr != nil {
		res.Data = "User Exists"
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, fmt.Errorf("el correo electronico ya esta siendo usado por otro usuario")
	}

	res.Data = "User no Exists"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h HandlerUsers) ValidateNickname(ctx context.Context, request *users_proto.ValidateNicknameRequest) (*users_proto.ValidateResponse, error) {
	res := &users_proto.ValidateResponse{Error: true}

	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	usrTemp, code, err := srvAuth.SrvUserTemp.GetUserByNickname(request.Nickname)
	if err != nil {
		logger.Error.Printf("couldn't get user by nickname: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	usr, code, err := srvAuth.SrvUser.GetUsersByNickname(request.Nickname)
	if err != nil {
		logger.Error.Printf("couldn't get user by nickname: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	if usrTemp != nil || usr != nil {
		res.Data = "User Exists"
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, fmt.Errorf("el nombre de usuario ya esta siendo usado por otro usuario")
	}

	res.Data = "User no Exists"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h HandlerUsers) GetUserById(ctx context.Context, request *users_proto.GetUserByIDRequest) (*users_proto.UserResponse, error) {
	res := &users_proto.UserResponse{Error: true}
	srvUser := auth.NewServerAuth(h.DB, nil, h.TxID)

	usr, code, err := srvUser.SrvUser.GetUsersByID(request.Id)
	if err != nil {
		logger.Error.Printf("couldn't get User by ID: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	if usr.FullPathPhoto != "" {
		pathPhoto := strings.Split(usr.FullPathPhoto, "-/")

		file, _, err := srvUser.SrvFiles.GetFileByPath(pathPhoto[0], pathPhoto[1])
		if err != nil {
			logger.Error.Printf("couldn't get profile picture: %v", err)
		}

		if file != nil {
			usr.FullPathPhoto = file.Encoding
		}
	}

	user := &users_proto.User{
		ID:         usr.ID,
		Nickname:   usr.Nickname,
		Email:      usr.Email,
		Name:       usr.Name,
		Lastname:   usr.Lastname,
		IdType:     int32(usr.IdType),
		IdNumber:   usr.IdNumber,
		Cellphone:  usr.Cellphone,
		StatusId:   int32(usr.StatusId),
		IdRole:     int32(usr.IdRole),
		RsaPrivate: usr.RsaPrivate,
		RsaPublic:  usr.RsaPublic,
	}

	res.Data = user
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, err
}

func (h HandlerUsers) ValidateIdentity(ctx context.Context, request *users_proto.ValidateIdentityRequest) (*users_proto.ValidateResponse, error) {
	res := &users_proto.ValidateResponse{Error: true}

	// TODO implements GetUserFromToken.
	userID := ""
	realIP := "c.IP()"

	srv := auth.NewServerAuth(h.DB, nil, h.TxID)
	// TODO implements GetUserFromToken.
	user, code, err := srv.SrvUser.GetUsersByID("u.ID")
	if err != nil {
		logger.Error.Printf("couldn't get user by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	if user == nil {
		logger.Error.Printf("couldn't get user by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	userWallet, code, err := srv.SrvUsersWallet.GetUserWalletByUserIDAndIdentityNumber(userID, request.IdentityNumber)
	if err != nil {
		logger.Error.Printf("couldn't get user wallet by user id and identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	if userWallet != nil {
		res.Data = "the user has already been verified"
		res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
		res.Error = false
		return res, err
	}

	_, code, err = srv.SrvUser.UpdateUsers(user.ID, user.Nickname, user.Email, user.Password, user.Name, user.Lastname, user.IdType, request.IdentityNumber, user.Cellphone, user.BirthDate, user.VerifiedCode, user.VerifiedAt, user.FullPathPhoto, user.RsaPrivate, user.RsaPublic, 21)

	if err != nil {
		logger.Error.Printf("couldn't update user by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	wallet, code, err := srv.SrvWallet.GetWalletByIdentityNumber(request.IdentityNumber)
	if err != nil {
		logger.Error.Printf("couldn't get wallet by identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	walletID := wallet.ID

	if wallet == nil {
		rsaPrivate, rsaPublic, _ := rsa_generate.Execute()
		rsaPrivateDevice, rsaPublicDevice, _ := rsa_generate.Execute()
		newWallet, code, err := srv.SrvWallet.CreateWallet(uuid.New().String(), mnemonic.Generate(), rsaPublic, rsaPrivate,
			rsaPublicDevice, rsaPrivateDevice, realIP, request.IdentityNumber, 1)
		if err != nil {
			logger.Error.Printf("couldn't create wallet: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return res, err
		}

		walletID = newWallet.ID

		_, code, err = srv.SrvAccounting.CreateAccounting(uuid.New().String(), newWallet.ID, 0, userID)
		if err != nil {
			logger.Error.Printf("couldn't create account to wallet: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return res, err
		}

	}

	_, code, err = srv.SrvUsersWallet.CreateUsersWallet(uuid.New().String(), userID, walletID, false)
	if err != nil {
		logger.Error.Printf("couldn't create users wallet: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	res.Data = "user verified"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h HandlerUsers) ValidateCertifier(ctx context.Context, request *users_proto.ValidateCertifierRequest) (*users_proto.ValidateResponse, error) {
	res := &users_proto.ValidateResponse{Error: true}

	// TODO implements GetUserFromToken.
	//u := helpers.GetUserContext(c)
	userID := ""
	realIP := ""

	srv := auth.NewServerAuth(h.DB, nil, h.TxID)
	user, code, err := srv.SrvUser.GetUsersByID(userID)
	if err != nil {
		logger.Error.Printf("couldn't get user by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	if user == nil {
		logger.Error.Printf("couldn't get user by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	userWallet, code, err := srv.SrvUsersWallet.GetUserWalletByUserIDAndIdentityNumber(userID, request.IdentityNumber)
	if err != nil {
		logger.Error.Printf("couldn't get user wallet by user id and identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	if userWallet != nil {
		res.Data = "the certifier has already been verified"
		res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
		res.Error = false
		return res, err
	}

	_, code, err = srv.SrvUser.UpdateUsers(user.ID, user.Nickname, user.Email, user.Password, user.Name, user.Lastname, user.IdType, request.IdentityNumber, user.Cellphone, user.BirthDate, user.VerifiedCode, user.VerifiedAt, user.FullPathPhoto, user.RsaPrivate, user.RsaPublic, 22)

	if err != nil {
		logger.Error.Printf("couldn't update user by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	wallet, code, err := srv.SrvWallet.GetWalletByIdentityNumber(request.IdentityNumber)
	if err != nil {
		logger.Error.Printf("couldn't get wallet by identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	walletID := wallet.ID

	if wallet == nil {
		rsaPrivate, rsaPublic, _ := rsa_generate.Execute()
		rsaPrivateDevice, rsaPublicDevice, _ := rsa_generate.Execute()
		newWallet, code, err := srv.SrvWallet.CreateWallet(uuid.New().String(), mnemonic.Generate(), rsaPublic, rsaPrivate,
			rsaPublicDevice, rsaPrivateDevice, realIP, request.IdentityNumber, 1)
		if err != nil {
			logger.Error.Printf("couldn't create wallet: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return res, err
		}

		walletID = newWallet.ID

		_, code, err = srv.SrvAccounting.CreateAccounting(uuid.New().String(), newWallet.ID, 0, userID)
		if err != nil {
			logger.Error.Printf("couldn't create account to wallet: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return res, err
		}

	}

	_, code, err = srv.SrvUsersWallet.CreateUsersWallet(uuid.New().String(), userID, walletID, false)
	if err != nil {
		logger.Error.Printf("couldn't create users wallet: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	res.Data = "certifier verified"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerUsers) ValidIdentityNumber(ctx context.Context, request *users_proto.RequestGetByIdentityNumber) (*users_proto.ResponseGetByIdentityNumber, error) {
	res := &users_proto.ResponseGetByIdentityNumber{Error: true}
	srvUser := auth.NewServerAuth(h.DB, nil, h.TxID)

	usrTemp, code, err := srvUser.SrvUserTemp.GetUserByIdentityNumber(request.IdentityNumber)
	if err != nil {
		logger.Error.Printf("couldn't get user by identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	usr, code, err := srvUser.SrvUser.GetUserByIdentityNumber(request.IdentityNumber)
	if err != nil {
		logger.Error.Printf("couldn't get User by identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	if usrTemp != nil || usr != nil {
		res.Data = "User Exist"
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, fmt.Errorf("el número de identificación ya esta siendo usado por otro usuario")
	}

	res.Data = "User no Exists"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerUsers) UpdateUserPhoto(ctx context.Context, request *users_proto.RequestUpdateUserPhoto) (*users_proto.ResponseUpdateUserPhoto, error) {
	res := &users_proto.ResponseUpdateUserPhoto{Error: true}
	u, err := helpers.GetUserContext(ctx)
	if err != nil {
		logger.Error.Printf("couldn't get token user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return res, err
	}

	srvUser := auth.NewServerAuth(h.DB, nil, h.TxID)

	idNumber, _ := strconv.ParseInt(u.IdNumber, 10, 64)
	f, err := srvUser.SrvFiles.UploadFile(idNumber, request.FileName, request.FileEncode)
	if err != nil {
		logger.Error.Printf("couldn't upload file s3: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return res, err
	}

	user := users.User{
		ID:            u.ID,
		FullPathPhoto: f.Path + "-/" + f.FileName,
		UpdatedAt:     time.Now(),
	}

	err = srvUser.SrvUser.ChangePicturePhoto(user)
	if err != nil {
		logger.Error.Printf("couldn't update profile picture: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		res.Msg = err.Error()
		return res, err
	}

	res.Data = "Foto de perfil correctamente actualizada"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerUsers) GetUserPhoto(ctx context.Context, request *users_proto.RequestGetUserPhoto) (*users_proto.ResponseGetUserPhoto, error) {
	res := &users_proto.ResponseGetUserPhoto{Error: true, Data: ""}
	u, err := helpers.GetUserContext(ctx)
	if err != nil {
		logger.Error.Printf("couldn't get token user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return res, err
	}

	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	user, code, err := srvAuth.SrvUser.GetUsersByID(u.ID)
	if err != nil {
		logger.Error.Printf("couldn't get user by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		res.Msg = err.Error()
		return res, err
	}

	if user.FullPathPhoto != "" {
		pathPhoto := strings.Split(user.FullPathPhoto, "-/")
		file, code, err := srvAuth.SrvFiles.GetFileByPath(pathPhoto[0], pathPhoto[1])
		if err != nil {
			logger.Error.Printf("couldn't get profile picture: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return res, err
		}

		res.Data = file.Encoding
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerUsers) ChangePassword(ctx context.Context, request *users_proto.RequestChangePwd) (*users_proto.ResponseChangePwd, error) {
	res := &users_proto.ResponseChangePwd{Error: true, Data: ""}
	u, err := helpers.GetUserContext(ctx)
	if err != nil {
		logger.Error.Printf("couldn't get token user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return res, err
	}

	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	if request.NewPassword != request.ConfirmPassword {
		logger.Error.Printf("new password and confirm password are not the same")
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return res, fmt.Errorf("new password and confirm password are not the same")
	}

	user, code, err := srvAuth.SrvUser.GetUsersByID(u.ID)
	if err != nil {
		logger.Error.Printf("couldn't get user by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		res.Msg = err.Error()
		return res, err
	}

	if !pwd.Compare(user.ID, user.Password, request.OldPassword) {
		logger.Error.Printf("old password is not correct")
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return res, fmt.Errorf("old password is not correct")
	}

	err = srvAuth.SrvUser.UpdatePassword(user.ID, request.NewPassword)
	if err != nil {
		logger.Error.Printf("couldn't update password: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return res, err
	}

	res.Data = "password updated"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}
