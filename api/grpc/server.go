package grpc

import (
	"blion-auth/api/grpc/handlers/accounting"
	"blion-auth/api/grpc/handlers/wallets"
	"blion-auth/internal/grpc/accounting_proto"
	"blion-auth/internal/grpc/wallet_proto"
	"blion-auth/pkg/auth/interceptor"
	"fmt"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"log"
	"net"

	"github.com/fatih/color"
)

const (
	version     = "0.0.1"
	website     = "https://www.bjungle.net/"
	banner      = `Blion Auht`
	description = `Blion Auth - %s - Port: %s
by bjungle 
Version: %s
%s`
)

type server struct {
	listening string
	DB        *sqlx.DB
	TxID      string
}

func newServer(listening int, db *sqlx.DB, txID string) *server {
	return &server{fmt.Sprintf(":%d", listening), db, txID}
}

func (srv *server) Start() {
	color.Blue(banner)
	color.Cyan(fmt.Sprintf(description, "", srv.listening, version, website))

	lis, err := net.Listen("tcp", "0.0.0.0"+srv.listening)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatalf("Error faltal listener %v", err)
	}

	itr := interceptor.NewAuthInterceptor()
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(itr.Unary()),
		grpc.StreamInterceptor(itr.Stream()),
	}

	s := grpc.NewServer(serverOptions...)

	wallet_proto.RegisterWalletServicesWalletServer(s, &wallets.HandlerWallet{DB: srv.DB, TxID: srv.TxID})
	accounting_proto.RegisterAccountingServicesAccountingServer(s, &accounting.HandlerAccounting{DB: srv.DB, TxID: srv.TxID})

	err = s.Serve(lis)
	if err != nil {
		log.Fatal("Error fatal server", err)
	}

}
