package users_temp

import (
	"fmt"
	"time"

	"blion-auth/internal/logger"
	"blion-auth/internal/models"

	"github.com/asaskevich/govalidator"
)

type PortsServerUsersTemp interface {
	CreateUsersTemp(id string, nickname string, email string, password string, name string, lastname string, idType int, idNumber string, cellphone string, birthDate time.Time, verifiedCode string) (*UserTemp, int, error)
	UpdateUsersTemp(id string, nickname string, email string, password string, name string, lastname string, idType int, idNumber string, cellphone string, birthDate time.Time, verifiedCode string) (*UserTemp, int, error)
	DeleteUsersTemp(id string) (int, error)
	GetUsersTempByID(id string) (*UserTemp, int, error)
	GetAllUsersTemp() ([]*UserTemp, error)
	GetUserByEmail(email string) (*UserTemp, int, error)
	GetUserByNickname(nickname string) (*UserTemp, int, error)
	GetUserByIdentityNumber(identityNumber string) (*UserTemp, int, error)
}

type service struct {
	repository ServicesUsersTempRepository
	user       *models.User
	txID       string
}

func NewUsersTempService(repository ServicesUsersTempRepository, user *models.User, TxID string) PortsServerUsersTemp {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateUsersTemp(id string, nickname string, email string, password string, name string, lastname string, idType int, idNumber string, cellphone string, birthDate time.Time, verifiedCode string) (*UserTemp, int, error) {
	m := NewUsersTemp(id, nickname, email, password, name, lastname, idType, idNumber, cellphone, birthDate, verifiedCode)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create UsersTemp :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateUsersTemp(id string, nickname string, email string, password string, name string, lastname string, idType int, idNumber string, cellphone string, birthDate time.Time, verifiedCode string) (*UserTemp, int, error) {
	m := NewUsersTemp(id, nickname, email, password, name, lastname, idType, idNumber, cellphone, birthDate, verifiedCode)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update UsersTemp :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteUsersTemp(id string) (int, error) {
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

func (s *service) GetUsersTempByID(id string) (*UserTemp, int, error) {
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

func (s *service) GetAllUsersTemp() ([]*UserTemp, error) {
	return s.repository.getAll()
}

func (s *service) GetUserByEmail(email string) (*UserTemp, int, error) {
	if !govalidator.IsEmail(email) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("email isn't email"))
		return nil, 15, fmt.Errorf("email isn't email")
	}
	m, err := s.repository.getByEmail(email)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetUserByNickname(nickname string) (*UserTemp, int, error) {
	m, err := s.repository.getByNickname(nickname)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetUserByIdentityNumber(identityNumber string) (*UserTemp, int, error) {
	m, err := s.repository.getByIdentityNumber(identityNumber)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
