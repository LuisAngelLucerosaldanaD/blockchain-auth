# BLion Auth

## Compilación

* para sistemas basados en Linux

````bash
GOOS=linux  GOARCH=amd64 go build 
````

* para sistemas basados en Windows

````bash
GOOS=windows  GOARCH=amd64 go build 
````


## Instalación de Protoc
Para instalar protoc dar click en este link [protoc](https://github.com/protocolbuffers/protobuf/releases)

## <span style="color:yellow">Atención</span>

Para versiones anterioeres a protoc v1.27.1 se ejecuta este comando
````bash
protoc -I api/grpc/proto --go_out=internal/grpc api/grpc/proto/*.proto
````
Para versiones superiores a protoc v1.27.1 se ejecuta este comando
````bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative api/grpc/proto/*.proto
````
