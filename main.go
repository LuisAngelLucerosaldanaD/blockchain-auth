package main

import (
	"blion-auth/api/grpc"
	"blion-auth/internal/env"
	"os"
)

func main() {
	c := env.NewConfiguration()
	_ = os.Setenv("AWS_ACCESS_KEY_ID", c.Aws.AWSACCESSKEYID)
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", c.Aws.AWSSECRETACCESSKEY)
	_ = os.Setenv("AWS_DEFAULT_REGION", c.Aws.AWSDEFAULTREGION)
	grpc.Start(c.App.Port)
}
