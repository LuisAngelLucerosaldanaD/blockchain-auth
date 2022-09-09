package users

import (
	"database/sql"
	"fmt"
	"time"

	"blion-auth/internal/models"

	"github.com/jmoiron/sqlx"
)

// psql estructura de conexión a la BD de postgresql
type psql struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func newUsersPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *User) error {
	date := time.Now()
	m.UpdatedAt = date
	m.CreatedAt = date
	m.IdUser = m.ID
	const psqlInsert = `INSERT INTO auth.users (id ,nickname, email, password, name, lastname, id_type, id_number, cellphone, status_id, failed_attempts, block_date, disabled_date, last_login, last_change_password, birth_date, verified_code, verified_at, is_deleted, id_user, full_path_photo, rsa_public, rsa_private, recovery_account_at, deleted_at, created_at, updated_at, id_role) VALUES (:id ,:nickname, :email, :password, :name, :lastname, :id_type, :id_number, :cellphone, :status_id, :failed_attempts, :block_date, :disabled_date, :last_login, :last_change_password, :birth_date, :verified_code, :verified_at, :is_deleted, :id_user, :full_path_photo, :rsa_public, :rsa_private, :recovery_account_at, :deleted_at,:created_at, :updated_at, :id_role) `
	rs, err := s.DB.NamedExec(psqlInsert, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) update(m *User) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE auth.users SET nickname = :nickname, email = :email, password = :password, name = :name, lastname = :lastname, id_type = :id_type, id_number = :id_number, cellphone = :cellphone, status_id = :status_id, failed_attempts = :failed_attempts, block_date = :block_date, disabled_date = :disabled_date, last_login = :last_login, last_change_password = :last_change_password, birth_date = :birth_date, verified_code = :verified_code, verified_at = :verified_at,  is_deleted = :is_deleted, id_user = :id_user, full_path_photo = :full_path_photo, rsa_public = :rsa_public, rsa_private = :rsa_private, recovery_account_at = :recovery_account_at, deleted_at = :deleted_at, updated_at = :updated_at, id_role = :id_role WHERE id = :id `
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) updateIdentity(m *User) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE auth.users SET name = :name, lastname = :lastname, id_number = :id_number, id_role = :id_role WHERE id = :id`
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *psql) delete(id string) error {
	const psqlDelete = `DELETE FROM auth.users WHERE id = :id `
	m := User{ID: id}
	rs, err := s.DB.NamedExec(psqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *psql) getByID(id string) (*User, error) {
	const psqlGetByID = `SELECT id , nickname, email, password, name, lastname, id_type, id_number, cellphone, status_id, failed_attempts, block_date, disabled_date, last_login, last_change_password, birth_date, verified_code, verified_at, is_deleted, id_user, id_role, full_path_photo,  rsa_public, rsa_private, recovery_account_at, deleted_at, created_at, updated_at FROM auth.users WHERE id = $1 `
	mdl := User{}
	err := s.DB.Get(&mdl, psqlGetByID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *psql) getAll() ([]*User, error) {
	var ms []*User
	const psqlGetAll = ` SELECT id , nickname, email, password, name, lastname, id_type, id_number, cellphone, status_id, failed_attempts, block_date, disabled_date, last_login, last_change_password, birth_date, verified_code, verified_at, is_deleted, id_user, id_role, full_path_photo,  rsa_public, rsa_private, recovery_account_at, deleted_at, created_at, updated_at FROM auth.users `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *psql) getByEmail(email string) (*User, error) {
	const psqlGetByEmail = `SELECT id , nickname, email, password, name, lastname, id_type, id_number, cellphone, status_id, failed_attempts, block_date, disabled_date, last_login, last_change_password, birth_date, verified_code, verified_at, is_deleted, id_user, id_role, full_path_photo,  rsa_public, rsa_private, recovery_account_at, deleted_at, created_at, updated_at FROM auth.users WHERE email = $1 `
	mdl := User{}
	err := s.DB.Get(&mdl, psqlGetByEmail, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

func (s *psql) getByNickname(nickname string) (*User, error) {
	const psqlGetByNickname = `SELECT id , nickname, email, password, name, lastname, id_type, id_number, cellphone, status_id, failed_attempts, block_date, disabled_date, last_login, last_change_password, birth_date, verified_code, verified_at, is_deleted, id_user, id_role, full_path_photo,  rsa_public, rsa_private, recovery_account_at, deleted_at, created_at, updated_at FROM auth.users WHERE nickname = $1 `
	mdl := User{}
	err := s.DB.Get(&mdl, psqlGetByNickname, nickname)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

func (s *psql) getByIdentityNumber(identityNumber string) (*User, error) {
	const psqlGetByIdentityNumber = `SELECT id , nickname, email, password, name, lastname, id_type, id_number, cellphone, status_id, failed_attempts, block_date, disabled_date, last_login, last_change_password, birth_date, verified_code, verified_at, is_deleted, id_user, id_role, full_path_photo, rsa_private, rsa_public, recovery_account_at, deleted_at, created_at, updated_at FROM auth.users WHERE id_number = $1 `
	mdl := User{}
	err := s.DB.Get(&mdl, psqlGetByIdentityNumber, identityNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

func (s *psql) updateProfilePhoto(user User) error {
	const psqlUpdate = `UPDATE auth.users SET full_path_photo = :full_path_photo, updated_at = :updated_at WHERE id = :id `
	rs, err := s.DB.NamedExec(psqlUpdate, &user)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

func (s *psql) updatePassword(userId, password string) error {
	user := User{
		ID:        userId,
		Password:  password,
		UpdatedAt: time.Now(),
	}
	const psqlUpdate = `UPDATE auth.users SET password = :password, updated_at = :updated_at WHERE id = :id `
	rs, err := s.DB.NamedExec(psqlUpdate, &user)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}
