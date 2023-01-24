package wallets

import (
	"blion-auth/internal/ciphers"
	"blion-auth/internal/env"
	"blion-auth/internal/grpc/wallet_proto"
	"blion-auth/internal/helpers"
	"blion-auth/internal/logger"
	"blion-auth/internal/mnemonic"
	"blion-auth/internal/models"
	"blion-auth/internal/msg"
	"blion-auth/internal/password"
	"blion-auth/internal/send_grid"
	genTemplate "blion-auth/internal/template"
	"blion-auth/pkg/auth"
	"blion-auth/pkg/auth/login"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type HandlerWallet struct {
	DB   *sqlx.DB
	TxID string
}

func (h *HandlerWallet) CreateWallet(ctx context.Context, request *wallet_proto.RequestCreateWallet) (*wallet_proto.ResponseCreateWallet, error) {
	res := &wallet_proto.ResponseCreateWallet{Error: true}
	e := env.NewConfiguration()
	var mailAttachment []*mail.Attachment

	var parameters = make(map[string]string, 0)
	u, err := helpers.GetUserContext(ctx)
	if err != nil {
		logger.Error.Printf("couldn't get token user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}

	srv := auth.NewServerAuth(h.DB, u, h.TxID)
	id := uuid.New().String()

	crtWallet, code, err := srv.SrvWallet.GetWalletByUserID(u.ID)
	if err != nil {
		logger.Error.Printf("No se pudo obtener la wallet por id, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}

	if crtWallet != nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(10009, h.DB, h.TxID)
		return res, err
	}

	mnemonicStr := mnemonic.Generate()
	w, code, err := srv.SrvWalletTemp.CreateWalletTemp(id, mnemonicStr, u.ID)
	if err != nil {
		logger.Error.Printf("couldn't create wallets temp: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	var tos []send_grid.To
	tos = append(tos, send_grid.To{Mail: u.Email, Name: "guest"})

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

	rsaPrivate, rsaPublic, err := ciphers.GenerateKeyPairEcdsa()
	if err != nil {
		logger.Error.Printf("No se pudo generar las claves ECDSA: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, fmt.Errorf("no se pudo generar las claves ECDSA")
	}
	wallet, code, err := srvAuth.SrvWallet.CreateWallet(walletTemp.ID, walletTemp.Mnemonic, rsaPublic, "127.0.0.1", u.IdNumber, 1)
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

	resWallet := wallet_proto.WalletActive{
		Id:       wallet.ID,
		Mnemonic: wallet.Mnemonic,
		Key: &wallet_proto.KeyPair{
			Public:  rsaPublic,
			Private: rsaPrivate,
		},
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Data = &resWallet
	res.Error = false
	return res, nil
}

func (h *HandlerWallet) GetWalletById(ctx context.Context, request *wallet_proto.RequestGetWalletById) (*wallet_proto.ResponseGetWalletById, error) {
	res := &wallet_proto.ResponseGetWalletById{Error: true}
	if request.Id == "" {
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		res.Error = false
		return res, fmt.Errorf("id is required")
	}

	srv := auth.NewServerAuth(h.DB, nil, h.TxID)

	wt, code, err := srv.SrvWallet.GetWalletByID(request.Id)
	if err != nil {
		logger.Error.Printf("couldn't get wallets by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	if wt == nil {
		logger.Error.Printf("couldn't get wallets by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, fmt.Errorf("no se encontro una wallet con el id proporcionado")
	}

	wallet := wallet_proto.Wallet{
		Id:             wt.ID,
		Mnemonic:       wt.Mnemonic,
		Public:         wt.RsaPublic,
		IpDevice:       wt.IpDevice,
		StatusId:       int32(wt.StatusId),
		IdentityNumber: wt.IdentityNumber,
		CreatedAt:      wt.CreatedAt.String(),
		UpdatedAt:      wt.UpdatedAt.String(),
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

	if wt == nil {
		res.Error = false
		res.Code, res.Type, res.Msg = msg.GetByCode(10007, h.DB, h.TxID)
		return res, err
	}

	wallet := wallet_proto.Wallet{
		Id:             wt.ID,
		Mnemonic:       wt.Mnemonic,
		Public:         wt.RsaPublic,
		IpDevice:       wt.IpDevice,
		StatusId:       int32(wt.StatusId),
		IdentityNumber: wt.IdentityNumber,
		CreatedAt:      wt.CreatedAt.String(),
		UpdatedAt:      wt.UpdatedAt.String(),
	}

	res.Data = &wallet
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
		Id:             wt.ID,
		Mnemonic:       wt.Mnemonic,
		Public:         wt.RsaPublic,
		IpDevice:       wt.IpDevice,
		StatusId:       int32(wt.StatusId),
		IdentityNumber: wt.IdentityNumber,
		CreatedAt:      wt.CreatedAt.String(),
		UpdatedAt:      wt.UpdatedAt.String(),
	}

	res.Data = &wallet
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerWallet) CreateWalletBySystem(ctx context.Context, request *wallet_proto.RqCreateWalletBySystem) (*wallet_proto.ResponseCreateWalletBySystem, error) {
	res := &wallet_proto.ResponseCreateWalletBySystem{Error: true}
	srv := auth.NewServerAuth(h.DB, nil, h.TxID)

	rsaPrivate, rsaPublic, err := ciphers.GenerateKeyPairEcdsa()
	if err != nil {
		logger.Error.Printf("No se pudo generar las claves ECDSA: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, fmt.Errorf("no se pudo generar las claves ECDSA")
	}
	wallet, code, err := srv.SrvWallet.CreateWallet(uuid.New().String(), mnemonic.Generate(), rsaPublic, "127.0.0.1", request.IdentityNumber, 1)
	if err != nil {
		logger.Error.Printf("couldn't create wallet: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	res.Data = &wallet_proto.DataWallet{
		Id:       wallet.ID,
		Mnemonic: wallet.Mnemonic,
		Key: &wallet_proto.KeyPair{
			Public:  rsaPublic,
			Private: rsaPrivate,
		},
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerWallet) UpdateWallet(ctx context.Context, request *wallet_proto.RqUpdateWallet) (*wallet_proto.ResUpdateWallet, error) {
	res := &wallet_proto.ResUpdateWallet{Error: true}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	newMnemonic := mnemonic.Generate()
	wallet, code, err := srvAuth.SrvWallet.UpdateWallet(request.Id, newMnemonic, request.IpDevice, request.IdentityNumber, int(request.StatusId))
	if err != nil {
		logger.Error.Printf("couldn't get update wallet, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	res.Data = &wallet_proto.Wallet{
		Id:             wallet.ID,
		Mnemonic:       newMnemonic,
		Public:         wallet.RsaPublic,
		IpDevice:       wallet.IpDevice,
		StatusId:       int32(wallet.StatusId),
		IdentityNumber: wallet.IdentityNumber,
		CreatedAt:      wallet.CreatedAt.String(),
		UpdatedAt:      wallet.UpdatedAt.String(),
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerWallet) FrozenMoney(ctx context.Context, request *wallet_proto.RqFrozenMoney) (*wallet_proto.ResFrozenMoney, error) {
	res := &wallet_proto.ResFrozenMoney{Error: true}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	u, err := helpers.GetUserContext(ctx)
	if err != nil {
		logger.Error.Printf("couldn't get token user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return res, err
	}

	_, code, err := srvAuth.SrvFrozenMoney.CreateFrozenMoney(uuid.New().String(), request.WalletId, request.Amount, request.LotteryId)
	if err != nil {
		logger.Error.Printf("couldn't frozen money, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	account, code, err := srvAuth.SrvAccounting.GetAccountingByWalletID(request.WalletId)
	if err != nil {
		logger.Error.Printf("couldn't get account, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	_, code, err = srvAuth.SrvAccounting.SetAmount(request.WalletId, account.Amount-request.Amount, u.ID)
	if err != nil {
		logger.Error.Printf("couldn't set account, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerWallet) UnFreezeMoney(ctx context.Context, request *wallet_proto.RqUnFreezeMoney) (*wallet_proto.ResUnFreezeMoney, error) {
	res := &wallet_proto.ResUnFreezeMoney{Error: true}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	u, err := helpers.GetUserContext(ctx)
	if err != nil {
		logger.Error.Printf("couldn't get token user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return res, err
	}

	frozenMoney, code, err := srvAuth.SrvFrozenMoney.GetFrozenMoneyByWalletIDAndLotteryId(request.WalletId, request.LotteryId)
	if err != nil {
		logger.Error.Printf("couldn't frozen money, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	if frozenMoney == nil {
		logger.Error.Printf("No tiene dinero congelado ha devolver")
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	account, code, err := srvAuth.SrvAccounting.GetAccountingByWalletID(request.WalletId)
	if err != nil {
		logger.Error.Printf("couldn't get account, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	amount := account.Amount + frozenMoney.Amount

	if request.Penalty > 0 {
		amount -= request.Penalty
	}

	_, code, err = srvAuth.SrvAccounting.SetAmount(request.WalletId, amount, u.ID)
	if err != nil {
		logger.Error.Printf("couldn't set account, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	_, err = srvAuth.SrvFrozenMoney.DeleteFrozenMoney(frozenMoney.ID)
	if err != nil {
		logger.Error.Printf("couldn't delete frozen money, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerWallet) GetFrozenMoney(ctx context.Context, request *wallet_proto.RqGetFrozenMoney) (*wallet_proto.ResGetFrozenMoney, error) {
	res := &wallet_proto.ResGetFrozenMoney{Error: true, Data: 0}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	frozenMoney, code, err := srvAuth.SrvFrozenMoney.GetFrozenMoneyByWalletID(request.WalletId)
	if err != nil {
		logger.Error.Printf("couldn't frozen money, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	if frozenMoney != nil {
		res.Data = frozenMoney.Amount
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}
