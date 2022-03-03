package wallets_temp

import (
	"blion-auth/internal/password"
	"fmt"

	"blion-auth/internal/logger"
	"blion-auth/internal/models"

	"github.com/asaskevich/govalidator"
)

type PortsServerWalletTemp interface {
	CreateWalletTemp(id string, mnemonic string, idUser string) (*WalletTemp, int, error)
	UpdateWalletTemp(id string, mnemonic string, idUser string) (*WalletTemp, int, error)
	DeleteWalletTemp(id string) (int, error)
	GetWalletTempByID(id string) (*WalletTemp, int, error)
	GetAllWalletTemp() ([]*WalletTemp, error)
	GetWalletTempByUserID(id, idUser string) (*WalletTemp, int, error)
}

type service struct {
	repository ServicesWalletTempRepository
	user       *models.User
	txID       string
}

func NewWalletTempService(repository ServicesWalletTempRepository, user *models.User, TxID string) PortsServerWalletTemp {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateWalletTemp(id string, mnemonic string, idUser string) (*WalletTemp, int, error) {

	m := NewWalletTemp(id, mnemonic, idUser)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	m.Mnemonic = password.Encrypt(m.Mnemonic)
	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create WalletTemp :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateWalletTemp(id string, mnemonic string, idUser string) (*WalletTemp, int, error) {
	m := NewWalletTemp(id, mnemonic, idUser)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update WalletTemp :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteWalletTemp(id string) (int, error) {
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

func (s *service) GetWalletTempByID(id string) (*WalletTemp, int, error) {
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

func (s *service) GetAllWalletTemp() ([]*WalletTemp, error) {
	return s.repository.getAll()
}

func (s *service) GetWalletTempByUserID(id, idUser string) (*WalletTemp, int, error) {
	if !govalidator.IsUUID(idUser) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return nil, 15, fmt.Errorf("id isn't uuid")
	}
	m, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	if m == nil {
		return nil, 22, err
	}
	if m.IdUser != idUser {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
