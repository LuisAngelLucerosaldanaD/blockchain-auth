




protoc -I api/grpc/proto --go_out=plugins=grpc:internal/grpc api/grpc/proto/*.proto
