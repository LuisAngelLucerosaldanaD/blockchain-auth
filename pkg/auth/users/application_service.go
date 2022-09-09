package users

import (
	"blion-auth/internal/pwd"
	"fmt"
	"time"

	"blion-auth/internal/logger"
	"blion-auth/internal/models"

	"github.com/asaskevich/govalidator"
)

type PortsServerUsers interface {
	CreateUsers(id string, nickname string, email string, password string, name string, lastname string, idType int, idNumber string, cellphone string, birthDate time.Time, verifiedCode string, verifiedAt time.Time, fullPathPhoto string, rsaPrivate string, rsaPublic string, idRole int) (*User, int, error)
	UpdateUsers(id string, nickname string, email string, password string, name string, lastname string, idType int, idNumber string, cellphone string, birthDate time.Time, verifiedCode string, verifiedAt time.Time, fullPathPhoto string, rsaPrivate string, rsaPublic string, idRole int) (*User, int, error)
	DeleteUsers(id string) (int, error)
	GetUsersByID(id string) (*User, int, error)
	GetUsersByEmail(email string) (*User, int, error)
	GetUsersByNickname(nickname string) (*User, int, error)
	GetAllUsers() ([]*User, error)
	GetUserByIdentityNumber(identityNumber string) (*User, int, error)
	ChangePicturePhoto(user User) error
	UpdatePassword(userID, password string) error
	UpdateUserIdentity(name string, lastname string, idNumber string, idRole int) (*User, int, error)
}

type service struct {
	repository ServicesUsersRepository
	user       *models.User
	txID       string
}

func NewUsersService(repository ServicesUsersRepository, user *models.User, TxID string) PortsServerUsers {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateUsers(id string, nickname string, email string, password string, name string, lastname string, idType int, idNumber string, cellphone string, birthDate time.Time, verifiedCode string, verifiedAt time.Time, fullPathPhoto string, rsaPrivate string, rsaPublic string, idRole int) (*User, int, error) {
	m := NewUsers(id, nickname, email, password, name, lastname, idType, idNumber, cellphone, birthDate, verifiedCode, verifiedAt, fullPathPhoto, rsaPrivate, rsaPublic, idRole)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Users :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateUsers(id string, nickname string, email string, password string, name string, lastname string, idType int, idNumber string, cellphone string, birthDate time.Time, verifiedCode string, verifiedAt time.Time, fullPathPhoto string, rsaPrivate string, rsaPublic string, idRole int) (*User, int, error) {
	m := NewUsers(id, nickname, email, password, name, lastname, idType, idNumber, cellphone, birthDate, verifiedCode, verifiedAt, fullPathPhoto, rsaPrivate, rsaPublic, idRole)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Users :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) UpdateUserIdentity(name string, lastname string, idNumber string, idRole int) (*User, int, error) {
	m := UpdateUser(name, lastname, idNumber, idRole)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.updateIdentity(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Users :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteUsers(id string) (int, error) {
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

func (s *service) GetUsersByID(id string) (*User, int, error) {
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

func (s *service) GetAllUsers() ([]*User, error) {
	return s.repository.getAll()
}

func (s *service) GetUsersByEmail(email string) (*User, int, error) {
	if !govalidator.IsEmail(email) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("email isn't email"))
		return nil, 15, fmt.Errorf("email isn't email")
	}
	m, err := s.repository.getByEmail(email)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByEmail row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetUsersByNickname(nickname string) (*User, int, error) {
	m, err := s.repository.getByNickname(nickname)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByNickname row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) ValidateUserIde(nickname string) (*User, int, error) {
	m, err := s.repository.getByNickname(nickname)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByNickname row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetUserByIdentityNumber(identityNumber string) (*User, int, error) {
	m, err := s.repository.getByIdentityNumber(identityNumber)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t get user by identityNumber row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) ChangePicturePhoto(user User) error {
	return s.repository.updateProfilePhoto(user)
}

func (s *service) UpdatePassword(userID, password string) error {
	password = pwd.Encrypt(password)
	return s.repository.updatePassword(userID, password)
}
