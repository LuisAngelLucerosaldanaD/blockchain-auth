package wallets

import (
	"blion-auth/internal/ciphers"
	"blion-auth/internal/grpc/wallet_proto"
	"blion-auth/internal/helpers"
	"blion-auth/internal/logger"
	"blion-auth/internal/mnemonic"
	"blion-auth/internal/msg"
	"blion-auth/pkg/auth"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type HandlerWallet struct {
	DB   *sqlx.DB
	TxID string
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

func (h *HandlerWallet) GetWalletByIdentityNumber(ctx context.Context, request *wallet_proto.RqGetByIdentityNumber) (*wallet_proto.ResponseGetByIdentityNumber, error) {
	res := &wallet_proto.ResponseGetByIdentityNumber{Error: true}
	srv := auth.NewServerAuth(h.DB, nil, h.TxID)

	wt, code, err := srv.SrvWallet.GetWalletByIdentityNumber(request.IdentityNumber)
	if err != nil {
		logger.Error.Printf("couldn't get wallets by identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
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

func (h *HandlerWallet) CreateWallet(ctx context.Context, request *wallet_proto.RequestCreateWallet) (*wallet_proto.ResponseCreateWallet, error) {
	res := &wallet_proto.ResponseCreateWallet{Error: true}
	srv := auth.NewServerAuth(h.DB, nil, h.TxID)

	rsaPrivate, rsaPublic, err := ciphers.GenerateKeyPairEcdsaX25519()
	if err != nil {
		logger.Error.Printf("No se pudo generar las claves ECDSA: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, fmt.Errorf("no se pudo generar las claves ECDSA")
	}

	mnemonicData := mnemonic.Generate()

	// TODO validar el envio de ip
	wallet, code, err := srv.SrvWallet.CreateWallet(uuid.New().String(), mnemonicData, rsaPublic, "127.0.0.1", request.IdentityNumber, 1)
	if err != nil {
		logger.Error.Printf("couldn't create wallet: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	res.Data = &wallet_proto.DataWallet{
		Id:       wallet.ID,
		Mnemonic: mnemonicData,
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
