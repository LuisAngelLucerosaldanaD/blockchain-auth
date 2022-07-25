package users_temp

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

func newUsersTempPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *UserTemp) error {
	date := time.Now()
	m.UpdatedAt = date
	m.CreatedAt = date
	m.IdUser = m.ID
	const psqlInsert = `INSERT INTO auth.users_temp (id ,nickname, email, password, name, lastname, id_type, id_number, cellphone, birth_date, verified_code, is_deleted, id_user, deleted_at, created_at, updated_at) VALUES (:id ,:nickname, :email, :password, :name, :lastname, :id_type, :id_number, :cellphone, :birth_date, :verified_code, :is_deleted, :id_user, :deleted_at,:created_at, :updated_at) `
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
func (s *psql) update(m *UserTemp) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE auth.users_temp SET nickname = :nickname, email = :email, password = :password, name = :name, lastname = :lastname, id_type = :id_type, id_number = :id_number, cellphone = :cellphone, birth_date = :birth_date, verified_code = :verified_code, is_deleted = :is_deleted, id_user = :id_user, deleted_at = :deleted_at, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM auth.users_temp WHERE id = :id `
	m := UserTemp{ID: id}
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
func (s *psql) getByID(id string) (*UserTemp, error) {
	const psqlGetByID = `SELECT id , nickname, email, password, name, lastname, id_type, id_number, cellphone, birth_date, verified_code, is_deleted, id_user, deleted_at, created_at, updated_at FROM auth.users_temp WHERE id = $1 `
	mdl := UserTemp{}
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
func (s *psql) getAll() ([]*UserTemp, error) {
	var ms []*UserTemp
	const psqlGetAll = ` SELECT id , nickname, email, password, name, lastname, id_type, id_number, cellphone, birth_date, verified_code, is_deleted, id_user, deleted_at, created_at, updated_at FROM auth.users_temp `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *psql) getByEmail(email string) (*UserTemp, error) {
	const psqlGetByEmail = `SELECT id , nickname, email, password, name, lastname, id_type, id_number, cellphone, birth_date, verified_code, is_deleted, id_user, deleted_at, created_at, updated_at FROM auth.users_temp WHERE email = $1 `
	mdl := UserTemp{}
	err := s.DB.Get(&mdl, psqlGetByEmail, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

func (s *psql) getByNickname(nickname string) (*UserTemp, error) {
	const psqlGetByEmail = `SELECT id , nickname, email, password, name, lastname, id_type, id_number, cellphone, birth_date, verified_code, is_deleted, id_user, deleted_at, created_at, updated_at FROM auth.users_temp WHERE nickname = $1 `
	mdl := UserTemp{}
	err := s.DB.Get(&mdl, psqlGetByEmail, nickname)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

func (s *psql) getByIdentityNumber(identityNumber string) (*UserTemp, error) {
	const psqlGetByEmail = `SELECT id , nickname, email, password, name, lastname, id_type, id_number, cellphone, birth_date, verified_code, is_deleted, id_user, deleted_at, created_at, updated_at FROM auth.users_temp WHERE id_number = $1 `
	mdl := UserTemp{}
	err := s.DB.Get(&mdl, psqlGetByEmail, identityNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}
