package login

import (
	"blion-auth/api/grpc/handlers/login"
	"blion-auth/internal/dbx"
	pb "blion-auth/internal/grpc/auth_proto"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {

	dbx := dbx.GetConnection()
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterAuthServicesUsersServer(s, &login.HandlerLogin{DB: dbx, TxID: uuid.New().String()})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	nickName := "nexum"

	client := pb.NewAuthServicesUsersClient(conn)
	resp, err := client.Login(ctx, &pb.LoginRequest{
		Nickname: &nickName,
		Password: "Nexum123",
	})
	if err != nil {
		t.Fatalf("Get block by id failed, error: %v", err)
	}
	log.Printf("Response: %+v", resp)
}
