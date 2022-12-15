# BLion Auth

## <span style="color:yellow">Atención</span>

Para versiones anterioeres a protoc v1.27.1 se ejecuta este comando
````bash
protoc -I api/grpc/proto --go_out=internal/grpc api/grpc/proto/*.proto
````
Para versiones superiores a protoc v1.27.1 se ejecuta este comando
````bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative api/grpc/proto/*.proto
````
