package grpc

import (
	"blion-auth/internal/dbx"
	"github.com/google/uuid"
)

func Start(port int) {
	db := dbx.GetConnection()
	defer db.Close()

	server := newServer(port, db, uuid.New().String())
	server.Start()
}
