#!/bin/bash
# Fallar en caso de error
set -e

echo "Generando código Protobuf..."
export PATH=$PATH:$(go env GOPATH)/bin

# Determinar raíz del repositorio relativo a este script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Rutas absolutas
PROTO_DIR="$REPO_ROOT/messaging/proto"
GW_PB_DIR="$REPO_ROOT/apps/gateway-service/pb"
WRITER_PB_DIR="$REPO_ROOT/apps/event-writer/pb"

# Crear directorios destino dentro de la estructura `apps/`
mkdir -p "$GW_PB_DIR"
mkdir -p "$WRITER_PB_DIR"

# Verificar existencia de .proto
if [ ! -f "$PROTO_DIR/blackfriday.proto" ]; then
  echo "Error: $PROTO_DIR/blackfriday.proto no existe"
  exit 1
fi

# Generar código Go para el Cliente (gateway-service)
protoc --go_out="$GW_PB_DIR" --go_opt=paths=source_relative \
    --go-grpc_out="$GW_PB_DIR" --go-grpc_opt=paths=source_relative \
    --proto_path="$PROTO_DIR" "$PROTO_DIR/blackfriday.proto"

# Generar código Go para el Servidor (event-writer)
protoc --go_out="$WRITER_PB_DIR" --go_opt=paths=source_relative \
    --go-grpc_out="$WRITER_PB_DIR" --go-grpc_opt=paths=source_relative \
    --proto_path="$PROTO_DIR" "$PROTO_DIR/blackfriday.proto"

echo "Código gRPC generado correctamente en:"
echo "  - apps/gateway-service/pb/blackfriday.pb.go"
echo "  - apps/event-writer/pb/blackfriday.pb.go"
