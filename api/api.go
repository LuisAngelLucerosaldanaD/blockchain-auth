package api

import "blion-auth/internal/dbx"

func Start(port int, app string, loggerHttp bool, allowedOrigins string) {
	db := dbx.GetConnection()
	defer db.Close()

	server := newServer(port)
	server.Start()
}
