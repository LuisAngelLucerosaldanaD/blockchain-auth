package users

import (
	"blion-auth/internal/env"
	"blion-auth/internal/grpc/users_proto"
	"blion-auth/internal/helpers"
	"blion-auth/internal/logger"
	"blion-auth/internal/models"
	"blion-auth/internal/msg"
	"blion-auth/internal/password"
	"blion-auth/internal/pwd"
	"blion-auth/internal/send_grid"
	genTemplate "blion-auth/internal/template"
	"blion-auth/pkg/auth"
	"blion-auth/pkg/auth/login"
	"blion-auth/pkg/auth/users"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type HandlerUsers struct {
	DB   *sqlx.DB
	TxID string
}

func (h *HandlerUsers) CreateUser(ctx context.Context, request *users_proto.UserRequest) (*users_proto.ResponseCreateUser, error) {
	res := &users_proto.ResponseCreateUser{Error: true}

	var parameters = make(map[string]string, 0)
	var mailAttachment []*mail.Attachment
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
	parameters["@access-code"] = emailCode
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

	var tos []send_grid.To
	tos = append(tos, send_grid.To{Mail: request.Email, Name: "guest"})

	usr.Password = ""
	usr.VerifiedCode = ""

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

	bodyCode, err := genTemplate.GenerateTemplateMail(parameters, e.Template.EmailToken)
	if err != nil {
		logger.Error.Printf("couldn't generate body in notification email")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}

	emailApi := send_grid.Model{
		FromMail:    e.SendGrid.FromMail,
		FromName:    e.SendGrid.FromName,
		Tos:         tos,
		Subject:     "Verificación de cuenta BLion",
		HTMLContent: bodyCode,
		Attachments: mailAttachment,
	}

	err = emailApi.SendMail()
	if err != nil {
		logger.Error.Println("error when execute send email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}

	request.Id = id
	request.Password = ""
	request.ConfirmPassword = ""

	res.Data = "Usuario creado correctamente, verifique su cuenta para poder iniciar sesión"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerUsers) ActivateUser(ctx context.Context, request *users_proto.ActivateUserRequest) (*users_proto.ValidateResponse, error) {
	res := &users_proto.ValidateResponse{Error: true}
	u, err := helpers.GetUserContext(ctx)
	if err != nil {
		logger.Error.Printf("couldn't get token user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
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

	verifiedAt := time.Now()
	usr, code, err := srvAuth.SrvUser.CreateUsers(usrTemp.ID, usrTemp.Nickname, usrTemp.Email, usrTemp.Password,
		usrTemp.Name, usrTemp.Lastname, usrTemp.IdType, usrTemp.IdNumber, usrTemp.Cellphone, usrTemp.BirthDate,
		"", &verifiedAt, "", 21)
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

	res.Data = "Active Account Successfully"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerUsers) ValidateEmail(ctx context.Context, request *users_proto.ValidateEmailRequest) (*users_proto.ValidateResponse, error) {
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
		res.Data = "el correo electronico ya esta siendo usado por otro usuario"
		res.Code, res.Type, res.Msg = msg.GetByCode(5, h.DB, h.TxID)
		return res, nil
	}

	res.Data = "User no Exists"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerUsers) ValidateNickname(ctx context.Context, request *users_proto.ValidateNicknameRequest) (*users_proto.ValidateResponse, error) {
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
		res.Data = "el nombre de usuario ya esta siendo usado por otro usuario"
		res.Code, res.Type, res.Msg = msg.GetByCode(5, h.DB, h.TxID)
		return res, nil
	}

	res.Data = "User no Exists"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerUsers) GetUserById(ctx context.Context, request *users_proto.GetUserByIDRequest) (*users_proto.UserResponse, error) {
	res := &users_proto.UserResponse{Error: true}
	srvUser := auth.NewServerAuth(h.DB, nil, h.TxID)

	usr, code, err := srvUser.SrvUser.GetUsersByID(request.Id)
	if err != nil {
		logger.Error.Printf("couldn't get User by ID: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	user := &users_proto.User{
		Id:            usr.ID,
		Nickname:      usr.Nickname,
		Email:         usr.Email,
		Name:          usr.Name,
		Lastname:      usr.Lastname,
		IdType:        int32(usr.IdType),
		IdNumber:      usr.IdNumber,
		Cellphone:     usr.Cellphone,
		StatusId:      int32(usr.StatusId),
		IdRole:        int32(usr.IdRole),
		FullPathPhoto: usr.FullPathPhoto,
	}

	res.Data = user
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, err
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
		res.Data = "el número de identificación ya esta siendo usado por otro usuario"
		res.Code, res.Type, res.Msg = msg.GetByCode(5, h.DB, h.TxID)
		return res, nil
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

func (h *HandlerUsers) CreateUserBySystem(ctx context.Context, request *users_proto.RequestCreateUserBySystem) (*users_proto.ResponseCreateUserBySystem, error) {
	res := &users_proto.ResponseCreateUserBySystem{Error: true}

	layout := "2006-01-02T15:04:05.000Z"
	var birthDate time.Time
	birthDate, err := time.Parse(layout, request.BirthDate)
	if err != nil {
		birthDate = time.Now()
	}

	verifiedAt := time.Now()
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	usr, code, err := srvAuth.SrvUser.CreateUsers(uuid.New().String(), request.Nickname, request.Email, password.Encrypt(request.Password),
		request.Name, request.Lastname, int(request.IdType), request.IdNumber, request.Cellphone, birthDate, "", &verifiedAt, "", 21)
	if err != nil {
		logger.Error.Printf("don't create user: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	res.Data = &users_proto.User{
		Id:        usr.ID,
		Nickname:  usr.Nickname,
		Email:     usr.Email,
		Name:      usr.Name,
		Lastname:  usr.Lastname,
		IdType:    int32(usr.IdType),
		IdNumber:  usr.IdNumber,
		Cellphone: usr.Cellphone,
		StatusId:  int32(usr.StatusId),
		BirthDate: usr.BirthDate.String(),
		IdRole:    int32(usr.IdRole),
	}
	res.Error = false
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	return res, nil
}

func (h *HandlerUsers) CreateUserWallet(ctx context.Context, request *users_proto.RqCreateUserWallet) (*users_proto.ResponseCreateUserWallet, error) {
	res := &users_proto.ResponseCreateUserWallet{Error: true}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	userWallet, code, err := srvAuth.SrvUsersWallet.CreateUsersWallet(uuid.New().String(), request.UserId, request.WalletId, false)
	if err != nil {
		logger.Error.Printf("don't create user wallet: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	res.Data = &users_proto.UserWallet{
		Id:        userWallet.ID,
		IdUser:    userWallet.IdUser,
		IdWallet:  userWallet.IdUser,
		IsDelete:  userWallet.IsDelete,
		DeletedAt: userWallet.DeletedAt.String(),
		CreatedAt: userWallet.CreatedAt.String(),
		UpdatedAt: userWallet.UpdatedAt.String(),
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerUsers) GetUserWalletByIdentityNumber(ctx context.Context, request *users_proto.RqGetUserWalletByIdentityNumber) (*users_proto.ResGetUserWalletByIdentityNumber, error) {
	res := &users_proto.ResGetUserWalletByIdentityNumber{Error: true}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	userWallet, code, err := srvAuth.SrvUsersWallet.GetUserWalletByUserIDAndIdentityNumber(request.UserId, request.IdentityNumber)
	if err != nil {
		logger.Error.Printf("couldn't get user wallet: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	if userWallet != nil {
		res.Data = &users_proto.UserWallet{
			Id:        userWallet.ID,
			IdUser:    userWallet.IdUser,
			IdWallet:  userWallet.IdUser,
			IsDelete:  userWallet.IsDelete,
			DeletedAt: userWallet.DeletedAt.String(),
			CreatedAt: userWallet.CreatedAt.String(),
			UpdatedAt: userWallet.UpdatedAt.String(),
		}
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerUsers) GetUserByIdentityNumber(ctx context.Context, request *users_proto.RqGetUserByIdentityNumber) (*users_proto.ResGetUserByIdentityNumber, error) {
	res := &users_proto.ResGetUserByIdentityNumber{Error: true}
	srvUser := auth.NewServerAuth(h.DB, nil, h.TxID)

	usr, code, err := srvUser.SrvUser.GetUserByIdentityNumber(request.IdentityNumber)
	if err != nil {
		logger.Error.Printf("couldn't get User by identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	user := &users_proto.User{
		Id:        usr.ID,
		Nickname:  usr.Nickname,
		Email:     usr.Email,
		Name:      usr.Name,
		Lastname:  usr.Lastname,
		IdType:    int32(usr.IdType),
		IdNumber:  usr.IdNumber,
		Cellphone: usr.Cellphone,
		StatusId:  int32(usr.StatusId),
		IdRole:    int32(usr.IdRole),
	}

	res.Data = user
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, err
}

func (h *HandlerUsers) UpdateUser(ctx context.Context, request *users_proto.RqUpdateUser) (*users_proto.ResUpdateUser, error) {
	res := &users_proto.ResUpdateUser{Error: true}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	layout := "2006-01-02T15:04:05.000Z"
	birthDate, err := time.Parse(layout, request.BirthDate)
	if err != nil {
		birthDate = time.Now()
	}

	verifiedAt, err := time.Parse(layout, request.VerifiedAt)
	if err != nil {
		verifiedAt = time.Now()
	}

	user, code, err := srvAuth.SrvUser.UpdateUsers(request.Id, request.Nickname, request.Email, request.Password, request.Name, request.Lastname, int(request.IdType), request.IdNumber, request.Cellphone, birthDate, request.VerifiedCode, &verifiedAt, request.FullPathPhoto, int(request.IdRole))
	if err != nil {
		logger.Error.Printf("couldn't update user: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	res.Data = &users_proto.User{
		Id:            user.ID,
		Nickname:      user.Nickname,
		Email:         user.Email,
		Name:          user.Name,
		Lastname:      user.Lastname,
		IdType:        int32(user.IdType),
		IdNumber:      user.IdUser,
		Cellphone:     user.Cellphone,
		StatusId:      int32(user.StatusId),
		BirthDate:     user.BirthDate.String(),
		IdUser:        user.IdUser,
		IdRole:        int32(user.IdRole),
		FullPathPhoto: user.FullPathPhoto,
		RealIp:        user.RealIP,
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerUsers) RequestChangePassword(ctx context.Context, request *users_proto.RqChangePwd) (*users_proto.ResAnny, error) {
	res := &users_proto.ResAnny{}
	e := env.NewConfiguration()
	var parameters = make(map[string]string, 0)
	var mailAttachment []*mail.Attachment

	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	usrI, code, err := srvAuth.SrvUser.GetUsersByEmail(request.Email)
	if err != nil {
		logger.Error.Printf("couldn't get user by identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return res, err
	}

	if usrI == nil {
		logger.Error.Printf("couldn't get user by identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, fmt.Errorf("no se pudo obtener el usuario por email")
	}

	var tos []send_grid.To
	tos = append(tos, send_grid.To{Mail: request.Email, Name: usrI.Name})

	usrI.Password = ""
	usrI.VerifiedCode = ""

	usrTemp := models.User{
		ID:        usrI.ID,
		Nickname:  usrI.Nickname,
		Email:     usrI.Email,
		Name:      usrI.Name,
		Lastname:  usrI.Lastname,
		IdType:    usrI.IdType,
		IdNumber:  usrI.IdNumber,
		Cellphone: usrI.Cellphone,
		BirthDate: usrI.BirthDate,
		CreatedAt: usrI.CreatedAt,
		UpdatedAt: usrI.UpdatedAt,
	}

	token, code, err := login.GenerateJWT(&usrTemp)
	if err != nil {
		logger.Error.Println("error, don't create token: %V", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	parameters["@url-token"] = e.Portal.Url + e.Portal.ChangePwd + token

	bodyCode, err := genTemplate.GenerateTemplateMail(parameters, e.Template.ChangePwd)
	if err != nil {
		logger.Error.Printf("couldn't generate body in notification email")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}

	emailApi := send_grid.Model{
		FromMail:    e.SendGrid.FromMail,
		FromName:    e.SendGrid.FromName,
		Tos:         tos,
		Subject:     "Verificación de cuenta BLion",
		HTMLContent: bodyCode,
		Attachments: mailAttachment,
	}

	err = emailApi.SendMail()
	if err != nil {
		logger.Error.Println("error when execute send email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}

	res.Data = "Se ha solicitado correctamente el cambio de contraseña"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerUsers) UpdateUserIdentity(ctx context.Context, request *users_proto.RqUpdateUserIdentity) (*users_proto.ResUpdateUser, error) {
	res := &users_proto.ResUpdateUser{Error: true}

	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)
	user, code, err := srvAuth.SrvUser.UpdateUserIdentity(request.Id, request.Name, request.Lastname, request.IdentityNumber, int(request.IdRole))
	if err != nil {
		logger.Error.Printf("couldn't update user: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	res.Data = &users_proto.User{
		Name:     user.Name,
		Lastname: user.Lastname,
		IdNumber: user.IdUser,
		IdRole:   int32(user.IdRole),
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}
