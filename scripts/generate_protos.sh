#!/bin/bash
# Fallar en caso de error
set -e

echo "Generando código Protobuf..."
export PATH=$PATH:$(go env GOPATH)/bin


# Crear directorios destino
mkdir -p go-grpc-client/pb
mkdir -p go-grpc-server/pb

# Generar código Go para el Cliente
# --go_out: genera structs y métodos básicos
# --go-grpc_out: genera interfaces de cliente y servidor gRPC
protoc --go_out=go-grpc-client/pb --go_opt=paths=source_relative \
    --go-grpc_out=go-grpc-client/pb --go-grpc_opt=paths=source_relative \
    --proto_path=proto proto/blackfriday.proto

# Generar código Go para el Servidor (mismo código, duplicado para independencia de módulos)
protoc --go_out=go-grpc-server/pb --go_opt=paths=source_relative \
    --go-grpc_out=go-grpc-server/pb --go-grpc_opt=paths=source_relative \
    --proto_path=proto proto/blackfriday.proto

echo " Código gRPC generado correctamente en:"
echo "   - go-grpc-client/pb/blackfriday.pb.go"
echo "   - go-grpc-server/pb/blackfriday.pb.go"
