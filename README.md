# Distributed Event-Driven Telemetry and Microservices Architecture

Academic project developed for the **Operating Systems 1** course at the Universidad de San Carlos de Guatemala (Second Semester 2025).

This repository contains a high-throughput, distributed event-driven telemetry and ingestion system built with Go microservices, gRPC, Apache Kafka, Redis, and Grafana for end-to-end observability.

---

## Overview

The system is designed to simulate, ingest, transport, and persist real-time telemetry events under heavy load. Ingestion requests are accepted via an HTTP gateway, dispatched over internal microservices via high-performance gRPC, streamed asynchronously through Apache Kafka message queues, cached using Redis, and monitored continuously with Grafana.

---

## Architecture

```
[ Client / Load Generator (Locust / k6) ]
                    |
                    v HTTP / REST
         [ Ingestion API Service ]
                    |
                    v gRPC (Protocol Buffers)
         [ Gateway / Router Service ]
                    |
                    v Event Stream
         [ Apache Kafka Message Broker ]
                    |
         +----------+----------+
         |                     |
         v                     v
[ Event Processor Worker ]  [ Event Writer Worker ]
         |                     |
         v                     v
    [ Redis Cache ]     [ Persistent Storage ]

Observability: Prometheus metrics collection with Grafana dashboards
Deployment: Docker Compose orchestration for local and containerized execution
```

---

## Key Features

- Low-latency internal inter-service communication using gRPC and Protocol Buffers
- Asynchronous decoupling and backpressure management via Apache Kafka partitions
- In-memory event state caching with Redis
- Comprehensive metrics instrumentation and Grafana dashboards for throughput and latency analysis
- Automated protobuf code generation pipelines
- Full Docker Compose environment definitions for reproducible deployments

---

## Technologies

- **Programming Language:** Go (Golang)
- **RPC Framework:** gRPC, Protocol Buffers
- **Messaging:** Apache Kafka
- **Caching:** Redis
- **Observability:** Prometheus, Grafana
- **Containerization:** Docker, Docker Compose

---

## Repository Structure

```
P2_SO1_202201336/
├── apps/
│   ├── ingestion-api/        # HTTP entrance service
│   ├── gateway-service/      # gRPC routing gateway
│   ├── event-processor/      # Real-time event consumer
│   └── event-writer/         # Persistence worker
├── deploy/                   # Docker Compose configurations
├── messaging/                # Kafka broker and topic configurations
├── observability/            # Grafana dashboard definitions
├── scripts/                  # Protobuf build and automation scripts
└── README.md
```

---

## Execution Guide

### Prerequisites
- Docker and Docker Compose
- Go 1.24+ (for local builds)
- Protocol Buffers compiler (`protoc`) with Go plugins

### Quick Start with Docker Compose

```bash
cd deploy/local-dev
docker compose -f docker-compose.yaml -f docker-compose.apps.yaml up --build
```

Access Grafana at `http://localhost:3000` to inspect real-time metrics.

---

## Academic Context

| Field       | Detail                                         |
|-------------|------------------------------------------------|
| Course      | Operating Systems 1 (SO1)                      |
| Institution | Universidad de San Carlos de Guatemala         |
| Semester    | Second Semester 2025                           |
| Author      | Estuardo Josue Vaquiax Reyes (202201336)       |
