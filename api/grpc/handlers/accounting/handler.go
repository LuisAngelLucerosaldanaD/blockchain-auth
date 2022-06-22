package accounting

import (
	"blion-auth/internal/grpc/accounting_proto"
	"blion-auth/internal/logger"
	"blion-auth/internal/msg"
	"blion-auth/pkg/auth"
	"context"
	"github.com/jmoiron/sqlx"
)

type HandlerAccounting struct {
	DB   *sqlx.DB
	TxID string
}

func (h HandlerAccounting) GetAccountingByWalletById(ctx context.Context, request *accounting_proto.RequestGetAccountingByWalletId) (*accounting_proto.ResponseGetAccountingByWalletId, error) {
	res := &accounting_proto.ResponseGetAccountingByWalletId{Error: true}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	wallet, code, err := srvAuth.SrvAccounting.GetAccountingByWalletID(request.Id)
	if err != nil {
		logger.Error.Printf("couldn't get accounting by wallet id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	res.Error = false
	res.Data = &accounting_proto.Accounting{
		Id:        wallet.ID,
		IdWallet:  wallet.IdWallet,
		Amount:    wallet.Amount,
		IdUser:    wallet.IdUser,
		CreatedAt: wallet.CreatedAt.String(),
		UpdatedAt: wallet.UpdatedAt.String(),
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)

	return res, nil
}

func (h HandlerAccounting) CreateAccounting(ctx context.Context, request *accounting_proto.RequestCreateAccounting) (*accounting_proto.ResponseCreateAccounting, error) {
	res := &accounting_proto.ResponseCreateAccounting{Error: true}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	wallet, code, err := srvAuth.SrvAccounting.CreateAccounting(request.Id, request.IdWallet, request.Amount, request.IdUser)
	if err != nil {
		logger.Error.Printf("couldn't create accounting: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	res.Error = false
	res.Data = &accounting_proto.Accounting{
		Id:        wallet.ID,
		IdWallet:  wallet.IdWallet,
		Amount:    wallet.Amount,
		IdUser:    wallet.IdUser,
		CreatedAt: wallet.CreatedAt.String(),
		UpdatedAt: wallet.UpdatedAt.String(),
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(30, h.DB, h.TxID)
	return res, nil
}

func (h HandlerAccounting) SetAmountToAccounting(ctx context.Context, request *accounting_proto.RequestSetAmountToAccounting) (*accounting_proto.ResponseSetAmountToAccounting, error) {
	res := &accounting_proto.ResponseSetAmountToAccounting{Error: true}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	wallet, code, err := srvAuth.SrvAccounting.SetAmount(request.WalletId, request.Amount, request.IdUser)
	if err != nil {
		logger.Error.Printf("couldn't set amount to accounting: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	res.Error = false
	res.Data = &accounting_proto.Accounting{
		Id:        wallet.ID,
		IdWallet:  wallet.IdWallet,
		Amount:    wallet.Amount,
		IdUser:    wallet.IdUser,
		CreatedAt: wallet.CreatedAt.String(),
		UpdatedAt: wallet.UpdatedAt.String(),
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(31, h.DB, h.TxID)
	return res, nil
}
