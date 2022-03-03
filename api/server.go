package api

import (
	"fmt"
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
}

func newServer(listening int) *server {
	return &server{fmt.Sprintf(":%d", listening)}
}

func (srv *server) Start() {
	color.Blue(banner)
	color.Cyan(fmt.Sprintf(description, "", srv.listening, version, website))

	lis, err := net.Listen("tcp", "0.0.0.0"+srv.listening)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatalf("Error faltal listener %v", err)
	}
	s := grpc.NewServer()

	//auth_proto.RegisterAuthServicesRolesServer(s, &roles.ServerRoles{})

	err = s.Serve(lis)
	if err != nil {
		log.Fatal("Error fatal server", err)
	}

}
