package main

import (
	"blion-auth/api/grpc"
	"blion-auth/internal/env"
)

func main() {
	c := env.NewConfiguration()
	grpc.Start(c.App.Port)
}
