package users_wallet

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

func newWalletPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *UsersWallet) error {
	date := time.Now()
	m.UpdatedAt = date
	m.CreatedAt = date
	const psqlInsert = `INSERT INTO auth.user_wallet (id, id_user, id_wallet, is_delete, deleted_at, created_at, updated_at) VALUES (:id ,:id_user, :id_wallet, :is_delete, :deleted_at, :created_at, :updated_at)`
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
func (s *psql) update(m *UsersWallet) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE auth.user_wallet SET id = :id, id_user = :id_user, id_wallet = :id_wallet, is_delete = :is_delete, deleted_at = :deleted_at, created_at = :created_at, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM auth.user_wallet WHERE id = :id `
	m := UsersWallet{ID: id}
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
func (s *psql) getByID(id string) (*UsersWallet, error) {
	const psqlGetByID = `SELECT id, id_user, id_wallet, is_delete, deleted_at, created_at, updated_at FROM auth.user_wallet WHERE id = $1 `
	mdl := UsersWallet{}
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
func (s *psql) getAll() ([]*UsersWallet, error) {
	var ms []*UsersWallet
	const psqlGetAll = `SELECT id, id_user, id_wallet, is_delete, deleted_at, created_at, updated_at FROM auth.user_wallet;`

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// getByUserIDAndIdentityNumber consulta un registro por su user id y identity number
func (s *psql) getByUserIDAndIdentityNumber(userID, identityNumber string) (*UsersWallet, error) {
	const psqlGetByUserIDAndIdentityNumber = `SELECT uw.id, uw.id_user, uw.id_wallet, uw.is_delete, uw.deleted_at, uw.created_at, uw.updated_at FROM auth.user_wallet uw join auth.wallets w on(uw.id_wallet = w.id) where w.identity_number = $1 and uw.id_user = $2;`
	mdl := UsersWallet{}
	err := s.DB.Get(&mdl, psqlGetByUserIDAndIdentityNumber, identityNumber, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}
