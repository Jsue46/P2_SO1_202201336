package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"

	pb "go-grpc-server/pb"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

type servidorGrpc struct {
	pb.UnimplementedProductSaleServiceServer
	escritorKafka *kafka.Writer
}

// ProcesarVenta recibe la petición gRPC y maneja la lógica de negocio
func (s *servidorGrpc) ProcesarVenta(ctx context.Context, req *pb.ProductSaleRequest) (*pb.ProductSaleResponse, error) {
	// Fase 2: Escribir en Kafka
	ventaBytes, err := json.Marshal(map[string]interface{}{
		"categoria":        req.Categoria.String(), // Convertir Enum a String (ej: "Electronica")
		"producto":         req.ProductoId,
		"precio":           req.Precio,
		"cantidad_vendida": req.CantidadVendida,
	})
	if err != nil {
		log.Printf("Error al serializar venta: %v", err)
		return nil, err
	}

	err = s.escritorKafka.WriteMessages(ctx, kafka.Message{
		Key:   []byte(req.ProductoId),
		Value: ventaBytes,
	})
	if err != nil {
		log.Printf("Error al escribir en Kafka: %v", err)
		return &pb.ProductSaleResponse{Estado: "Error al escribir en Kafka"}, nil
	}

	log.Printf("Venta enviada a Kafka: %s, %s", req.Categoria, req.ProductoId)
	return &pb.ProductSaleResponse{Estado: "Procesado y enviado a Cola"}, nil
}

func main() {
	addrKafka := os.Getenv("KAFKA_BROKER")
	if addrKafka == "" {
		addrKafka = "my-cluster-kafka-bootstrap:9092"
	}

	w := &kafka.Writer{
		Addr:     kafka.TCP(addrKafka),
		Topic:    "sales-topic",
		Balancer: &kafka.LeastBytes{},
	}
	defer w.Close()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error al escuchar puerto: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterProductSaleServiceServer(s, &servidorGrpc{escritorKafka: w})

	log.Printf("Servidor gRPC escuchando en %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
