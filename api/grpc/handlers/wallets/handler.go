package wallets

import (
	"blion-auth/internal/email"
	"blion-auth/internal/env"
	"blion-auth/internal/grpc/wallet_proto"
	"blion-auth/internal/helpers"
	"blion-auth/internal/logger"
	"blion-auth/internal/mnemonic"
	"blion-auth/internal/models"
	"blion-auth/internal/msg"
	"blion-auth/internal/password"
	"blion-auth/internal/rsa_generate"
	"blion-auth/pkg/auth"
	"blion-auth/pkg/auth/login"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type HandlerWallet struct {
	DB   *sqlx.DB
	TxID string
}

func (h *HandlerWallet) CreateWallet(ctx context.Context, request *wallet_proto.RequestCreateWallet) (*wallet_proto.ResponseCreateWallet, error) {
	res := &wallet_proto.ResponseCreateWallet{Error: true}
	e := env.NewConfiguration()

	var parameters = make(map[string]string, 0)
	u, err := helpers.GetUserContext(ctx)
	if err != nil {
		logger.Error.Printf("couldn't get token user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return res, err
	}

	srv := auth.NewServerAuth(h.DB, u, h.TxID)
	id := uuid.New().String()

	mnemonicStr := mnemonic.Generate()
	w, code, err := srv.SrvWalletTemp.CreateWalletTemp(id, mnemonicStr, u.ID)
	if err != nil {
		logger.Error.Printf("couldn't create wallets temp: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	tos := []string{u.Email}

	usrTemp := models.User{
		ID:        u.ID,
		Nickname:  u.Nickname,
		Email:     u.Email,
		Name:      u.Name,
		Lastname:  u.Lastname,
		IdType:    u.IdType,
		IdNumber:  u.IdNumber,
		Cellphone: u.Cellphone,
		BirthDate: u.BirthDate,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	token, code, err := login.GenerateJWT(&usrTemp)

	if err != nil {
		logger.Error.Println("error GenerateJWT: %V", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	parameters["@url-token"] = e.Portal.Url + e.Portal.ActivateWallet + token

	err = email.Send(tos, parameters, e.Template.EmailWalletToken, "Activate Wallet")
	if err != nil {
		logger.Error.Println("error when execute send email: %V", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(10002, h.DB, h.TxID)
		return res, err
	}

	resW := wallet_proto.DataWallet{
		Id:       w.ID,
		Mnemonic: mnemonicStr,
	}

	res.Data = &resW
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerWallet) ActivateWallet(ctx context.Context, request *wallet_proto.RequestActivateWallet) (*wallet_proto.ResponseActivateWallet, error) {
	res := &wallet_proto.ResponseActivateWallet{Error: true}
	u, err := helpers.GetUserContext(ctx)
	if err != nil {
		logger.Error.Printf("couldn't get token user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return res, err
	}

	var rsaPrivate, rsaPublic, rsaPrivateDevice, rsaPublicDevice string

	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	walletTemp, code, err := srvAuth.SrvWalletTemp.GetWalletTempByUserID(request.Id, u.ID)
	if err != nil {
		logger.Error.Printf("couldn't get wallet by UserId: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}
	if walletTemp == nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(10003, h.DB, h.TxID)
		return res, fmt.Errorf("wallet not found")
	}
	if !password.Compare(walletTemp.ID, walletTemp.Mnemonic, request.Mnemonic) {
		logger.Error.Printf("the verification mnemonic is not correct: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(10004, h.DB, h.TxID)
		return res, fmt.Errorf("the verification mnemonic is not correct")
	}
	//TODO Create worker RSA
	rsaPrivate, rsaPublic, err = rsa_generate.Execute()
	rsaPrivateDevice, rsaPublicDevice, err = rsa_generate.Execute()
	wallet, code, err := srvAuth.SrvWallet.CreateWallet(walletTemp.ID, walletTemp.Mnemonic, rsaPublic, rsaPrivate,
		rsaPublicDevice, rsaPrivateDevice, "127.0.0.1", u.IdNumber, 1)
	if err != nil {
		logger.Error.Printf("couldn't create wallet: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	code, err = srvAuth.SrvWalletTemp.DeleteWalletTemp(walletTemp.ID)
	if err != nil {
		logger.Error.Printf("couldn't delete wallet temp by UserId: %v", err)
	}

	_, code, err = srvAuth.SrvAccounting.CreateAccounting(uuid.New().String(), wallet.ID, 0, u.ID)
	if err != nil {
		logger.Error.Printf("couldn't create account to wallet: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	_, code, err = srvAuth.SrvUsersWallet.CreateUsersWallet(uuid.New().String(), u.ID, wallet.ID, false)
	if err != nil {
		logger.Error.Printf("couldn't create users wallet: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	resWallet := wallet_proto.Wallet{
		Id:               wallet.ID,
		Mnemonic:         wallet.Mnemonic,
		RsaPublic:        wallet.RsaPublic,
		RsaPrivate:       wallet.RsaPrivate,
		RsaPublicDevice:  wallet.RsaPublicDevice,
		RsaPrivateDevice: wallet.RsaPrivateDevice,
		IpDevice:         wallet.IpDevice,
		StatusId:         int32(wallet.StatusId),
		IdentityNumber:   wallet.IdentityNumber,
		CreatedAt:        wallet.CreatedAt.String(),
		UpdatedAt:        wallet.UpdatedAt.String(),
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Data = &resWallet
	res.Error = false
	return res, nil
}

func (h *HandlerWallet) GetWalletById(ctx context.Context, request *wallet_proto.RequestGetWalletById) (*wallet_proto.ResponseGetWalletById, error) {
	res := &wallet_proto.ResponseGetWalletById{Error: true}
	var id string
	if request.Id == "" {
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		res.Error = false
		return res, fmt.Errorf("id is required")
	}

	srv := auth.NewServerAuth(h.DB, nil, h.TxID)

	wt, _, err := srv.SrvWallet.GetWalletByID(id)
	if err != nil {
		logger.Error.Printf("couldn't get wallets by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return res, err
	}

	wallet := wallet_proto.Wallet{
		Id:               wt.ID,
		Mnemonic:         wt.Mnemonic,
		RsaPublic:        wt.RsaPublic,
		RsaPrivate:       wt.RsaPrivate,
		RsaPublicDevice:  wt.RsaPublicDevice,
		RsaPrivateDevice: wt.RsaPrivateDevice,
		IpDevice:         wt.IpDevice,
		StatusId:         int32(wt.StatusId),
		IdentityNumber:   wt.IdentityNumber,
		CreatedAt:        wt.CreatedAt.String(),
		UpdatedAt:        wt.UpdatedAt.String(),
	}

	res.Data = &wallet
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerWallet) GetWalletByUserId(ctx context.Context, request *wallet_proto.RequestGetWalletByUserId) (*wallet_proto.ResponseGetWalletByUserId, error) {
	res := &wallet_proto.ResponseGetWalletByUserId{Error: true}

	u, err := helpers.GetUserContext(ctx)
	if err != nil {
		logger.Error.Printf("couldn't get token user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return res, err
	}

	srv := auth.NewServerAuth(h.DB, nil, h.TxID)

	wt, _, err := srv.SrvWallet.GetWalletByUserID(u.ID)
	if err != nil {
		logger.Error.Printf("couldn't get wallets by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return res, err
	}

	var wallets []*wallet_proto.Wallet
	for _, wallet := range wt {
		wallets = append(wallets, &wallet_proto.Wallet{
			Id:               wallet.ID,
			Mnemonic:         wallet.Mnemonic,
			RsaPublic:        wallet.RsaPublic,
			RsaPrivate:       wallet.RsaPrivate,
			RsaPublicDevice:  wallet.RsaPublicDevice,
			RsaPrivateDevice: wallet.RsaPrivateDevice,
			IpDevice:         wallet.IpDevice,
			StatusId:         int32(wallet.StatusId),
			IdentityNumber:   wallet.IdentityNumber,
			CreatedAt:        wallet.CreatedAt.String(),
			UpdatedAt:        wallet.UpdatedAt.String(),
		})
	}

	res.Data = wallets
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerWallet) GetWalletByIdentityNumber(ctx context.Context, request *wallet_proto.RqGetByIdentityNumber) (*wallet_proto.ResponseGetByIdentityNumber, error) {
	res := &wallet_proto.ResponseGetByIdentityNumber{Error: true}
	srv := auth.NewServerAuth(h.DB, nil, h.TxID)

	wt, code, err := srv.SrvWallet.GetWalletByIdentityNumber(request.IdentityNumber)
	if err != nil {
		logger.Error.Printf("couldn't get wallets by identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	wallet := wallet_proto.Wallet{
		Id:               wt.ID,
		Mnemonic:         wt.Mnemonic,
		RsaPublic:        wt.RsaPublic,
		RsaPrivate:       wt.RsaPrivate,
		RsaPublicDevice:  wt.RsaPublicDevice,
		RsaPrivateDevice: wt.RsaPrivateDevice,
		IpDevice:         wt.IpDevice,
		StatusId:         int32(wt.StatusId),
		IdentityNumber:   wt.IdentityNumber,
		CreatedAt:        wt.CreatedAt.String(),
		UpdatedAt:        wt.UpdatedAt.String(),
	}

	res.Data = &wallet
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerWallet) CreateWalletBySystem(ctx context.Context, request *wallet_proto.RqCreateWalletBySystem) (*wallet_proto.ResponseCreateWalletBySystem, error) {
	res := &wallet_proto.ResponseCreateWalletBySystem{Error: true}
	srv := auth.NewServerAuth(h.DB, nil, h.TxID)

	rsaPrivate, rsaPublic, _ := rsa_generate.Execute()
	rsaPrivateDevice, rsaPublicDevice, _ := rsa_generate.Execute()
	wallet, code, err := srv.SrvWallet.CreateWallet(uuid.New().String(), mnemonic.Generate(), rsaPublic, rsaPrivate,
		rsaPublicDevice, rsaPrivateDevice, "127.0.0.1", request.IdentityNumber, 1)
	if err != nil {
		logger.Error.Printf("couldn't create wallet: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	res.Data = &wallet_proto.DataWallet{
		Id:       wallet.ID,
		Mnemonic: wallet.Mnemonic,
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}
