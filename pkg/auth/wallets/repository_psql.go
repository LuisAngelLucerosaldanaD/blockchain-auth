package wallets

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
func (s *psql) create(m *Wallet) error {
	date := time.Now()
	m.UpdatedAt = date
	m.CreatedAt = date
	const psqlInsert = `INSERT INTO auth.wallet (id ,mnemonic, rsa_public, ip_device, status_id, dni, created_at, updated_at) VALUES (:id ,:mnemonic, :rsa_public, :ip_device, :status_id, :dni, :created_at, :updated_at) `
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
func (s *psql) update(m *Wallet) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE auth.wallet SET id = :id ,mnemonic = :mnemonic, rsa_public = :rsa_public, ip_device = :ip_device, status_id = :status_id, dni = :dni, created_at = :created_at, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM auth.wallet WHERE id = :id `
	m := Wallet{ID: id}
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
func (s *psql) getByID(id string) (*Wallet, error) {
	const psqlGetByID = `SELECT id ,mnemonic, rsa_public, ip_device, status_id, dni, created_at, updated_at FROM auth.wallet WHERE id = $1 `
	mdl := Wallet{}
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
func (s *psql) getAll() ([]*Wallet, error) {
	var ms []*Wallet
	const psqlGetAll = ` SELECT id ,mnemonic, rsa_public, ip_device, status_id, dni, created_at, updated_at FROM auth.wallet `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// GetAll consulta todos los registros de la BD
func (s *psql) getWalletByIdentityNumber(identityNumber string) (*Wallet, error) {
	mdl := Wallet{}
	const psqlGetAll = `SELECT w.id, w.mnemonic, w.rsa_public, w.ip_device, w.status_id, w.dni, w.created_at, w.updated_at FROM auth.wallet w WHERE w.dni = $1 LIMIT 1;`

	err := s.DB.Get(&mdl, psqlGetAll, identityNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}
