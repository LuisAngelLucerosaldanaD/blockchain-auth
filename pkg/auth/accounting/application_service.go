package accounting

import (
	"fmt"

	"blion-auth/internal/logger"
	"blion-auth/internal/models"

	"github.com/asaskevich/govalidator"
)

type PortsServerAccounting interface {
	CreateAccounting(id string, idWallet string, amount float64) (*Accounting, int, error)
	UpdateAccounting(id string, idWallet string, amount float64) (*Accounting, int, error)
	DeleteAccounting(id string) (int, error)
	GetAccountingByID(id string) (*Accounting, int, error)
	GetAllAccounting() ([]*Accounting, error)
	SetAmount(idWallet string, amount float64) (*Accounting, int, error)
	GetAccountingByWalletID(walletID string) (*Accounting, int, error)
}

type service struct {
	repository ServicesAccountingRepository
	user       *models.User
	txID       string
}

func NewAccountingService(repository ServicesAccountingRepository, user *models.User, TxID string) PortsServerAccounting {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateAccounting(id string, idWallet string, amount float64) (*Accounting, int, error) {
	m := NewAccounting(id, idWallet, amount)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return nil, 15, fmt.Errorf("id isn't uuid")
	}
	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Accounting :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateAccounting(id string, idWallet string, amount float64) (*Accounting, int, error) {
	m := NewAccounting(id, idWallet, amount)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return nil, 15, fmt.Errorf("id isn't uuid")
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Accounting :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteAccounting(id string) (int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return 15, fmt.Errorf("id isn't uuid")
	}

	if err := s.repository.delete(id); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update row:", err)
		return 20, err
	}
	return 28, nil
}

func (s *service) GetAccountingByID(id string) (*Accounting, int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return nil, 15, fmt.Errorf("id isn't uuid")
	}
	m, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllAccounting() ([]*Accounting, error) {
	return s.repository.getAll()
}

func (s *service) SetAmount(idWallet string, amount float64) (*Accounting, int, error) {
	m := NewAccountingSetAmount(idWallet, amount)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.setAmount(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Accounting :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) GetAccountingByWalletID(walletID string) (*Accounting, int, error) {
	if !govalidator.IsUUID(walletID) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return nil, 15, fmt.Errorf("id isn't uuid")
	}
	m, err := s.repository.getByWalletID(walletID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
