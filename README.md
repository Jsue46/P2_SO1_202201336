# P2_SO1_202201336

Breve guía para preparar y ejecutar el proyecto localmente.

Requisitos:
- Docker y docker-compose
- Go 1.24 (para generar/compilar localmente)
- protoc + plugins `protoc-gen-go` y `protoc-gen-go-grpc` si quiere regenerar protos

Pasos rápidos:

1. Generar los protos Go (opcional si ya están en `apps/*/pb`):

```bash
cd deploy/local-dev
../scripts/generate_protos.sh
```

2. Levantar los servicios (modo local con Kafka/Redis y apps):

```bash
cd deploy/local-dev
docker compose -f docker-compose.yaml -f docker-compose.apps.yaml up --build
```

3. (Opcional) Levantar Grafana para observabilidad:

```bash
docker compose -f docker-compose.grafana.yaml up
```

Notas:
- Los servicios Go usan módulos locales (`go-grpc-server`, `go-grpc-client`, `go-consumer`) y los paquetes generados por `protoc` se encuentran en `apps/*/pb`.
- Si moviste archivos, asegúrate de regenerar protos con `scripts/generate_protos.sh`.
