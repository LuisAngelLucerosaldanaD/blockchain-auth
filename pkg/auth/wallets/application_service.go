package wallets

import (
	"blion-auth/internal/password"
	"fmt"

	"blion-auth/internal/logger"
	"blion-auth/internal/models"

	"github.com/asaskevich/govalidator"
)

type PortsServerWallet interface {
	CreateWallet(id, mnemonic, rsaPublic, rsaPrivate, rsaPublicDevice, rsaPrivateDevice, ipDevice, identityNumber string, statusId int) (*Wallet, int, error)
	UpdateWallet(id, mnemonic, rsaPublic, rsaPrivate, rsaPublicDevice, rsaPrivateDevice, ipDevice, identityNumber string, statusId int) (*Wallet, int, error)
	DeleteWallet(id string) (int, error)
	GetWalletByID(id string) (*Wallet, int, error)
	GetAllWallet() ([]*Wallet, error)
	GetWalletByUserID(userID string) ([]*Wallet, int, error)
	GetWalletByIdentityNumber(identityNumber string) (*Wallet, int, error)
}

type service struct {
	repository ServicesWalletRepository
	user       *models.User
	txID       string
}

func NewWalletService(repository ServicesWalletRepository, user *models.User, TxID string) PortsServerWallet {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateWallet(id, mnemonic, rsaPublic, rsaPrivate, rsaPublicDevice, rsaPrivateDevice, ipDevice, identityNumber string, statusId int) (*Wallet, int, error) {
	m := NewWallet(id, mnemonic, rsaPublic, rsaPrivate, rsaPublicDevice, rsaPrivateDevice, ipDevice, identityNumber, statusId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	m.Mnemonic = password.Encrypt(mnemonic)
	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Wallet :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateWallet(id, mnemonic, rsaPublic, rsaPrivate, rsaPublicDevice, rsaPrivateDevice, ipDevice, identityNumber string, statusId int) (*Wallet, int, error) {
	m := NewWallet(id, mnemonic, rsaPublic, rsaPrivate, rsaPublicDevice, rsaPrivateDevice, ipDevice, identityNumber, statusId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Wallet :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteWallet(id string) (int, error) {
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

func (s *service) GetWalletByID(id string) (*Wallet, int, error) {
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

func (s *service) GetAllWallet() ([]*Wallet, error) {
	return s.repository.getAll()
}

func (s *service) GetWalletByUserID(userID string) ([]*Wallet, int, error) {
	if !govalidator.IsUUID(userID) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return nil, 15, fmt.Errorf("id isn't uuid")
	}
	m, err := s.repository.getWalletByUserId(userID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetWalletByIdentityNumber(identityNumber string) (*Wallet, int, error) {
	if identityNumber == "" {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return nil, 15, fmt.Errorf("id is required")
	}
	m, err := s.repository.getWalletByIdentityNumber(identityNumber)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByIdentityNumber row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
