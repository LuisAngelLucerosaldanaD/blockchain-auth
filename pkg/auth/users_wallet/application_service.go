package users_wallet

import (
	"fmt"

	"blion-auth/internal/logger"
	"blion-auth/internal/models"

	"github.com/asaskevich/govalidator"
)

type PortsServerUsersWallet interface {
	CreateUsersWallet(id, idUser, idWallet string, isDelete bool) (*UsersWallet, int, error)
	UpdateUsersWallet(id, idUser, idWallet string, isDelete bool) (*UsersWallet, int, error)
	DeleteUsersWallet(id string) (int, error)
	GetUsersWalletByID(id string) (*UsersWallet, int, error)
	GetAllUsersWallet() ([]*UsersWallet, error)
	GetUserWalletByUserIDAndIdentityNumber(userID, identityNumber string) (*UsersWallet, int, error)
}

type service struct {
	repository ServicesUsersWalletRepository
	user       *models.User
	txID       string
}

func NewUsersWalletService(repository ServicesUsersWalletRepository, user *models.User, TxID string) PortsServerUsersWallet {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateUsersWallet(id, idUser, idWallet string, isDelete bool) (*UsersWallet, int, error) {
	m := NewWallet(id, idUser, idWallet, isDelete)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Wallet :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateUsersWallet(id, idUser, idWallet string, isDelete bool) (*UsersWallet, int, error) {
	m := NewWallet(id, idUser, idWallet, isDelete)
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

func (s *service) DeleteUsersWallet(id string) (int, error) {
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

func (s *service) GetUsersWalletByID(id string) (*UsersWallet, int, error) {
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

func (s *service) GetAllUsersWallet() ([]*UsersWallet, error) {
	return s.repository.getAll()
}

func (s *service) GetUserWalletByUserIDAndIdentityNumber(userID, identityNumber string) (*UsersWallet, int, error) {
	if !govalidator.IsUUID(userID) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return nil, 15, fmt.Errorf("id isn't uuid")
	}
	if identityNumber == "" {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("identityNumber required"))
		return nil, 15, fmt.Errorf("identityNumber required")
	}
	m, err := s.repository.getByUserIDAndIdentityNumber(userID, identityNumber)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByUserIDAndIdentityNumber row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
