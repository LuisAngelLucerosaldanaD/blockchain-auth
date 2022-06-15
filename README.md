# BLion Auth

Para generar los archivos .proto ejecutar el siguiente comando:
````bash
protoc -I api/grpc/proto --go_out=plugins=grpc:internal/grpc api/grpc/proto/*.proto
````
