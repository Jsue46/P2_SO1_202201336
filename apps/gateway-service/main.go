package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	pb "go-grpc-client/pb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SolicitudVenta struct {
	Categoria       string  `json:"categoria"`
	Producto        string  `json:"producto_id"`
	Precio          float64 `json:"precio"`
	CantidadVendida int32   `json:"cantidad_vendida"`
}

func main() {
	direccioGrpc := os.Getenv("GRPC_SERVER_ADDR")
	if direccioGrpc == "" {
		direccioGrpc = "localhost:50051"
	}

	// Conectar al servidor gRPC
	conn, err := grpc.NewClient(direccioGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("no se pudo conectar: %v", err)
	}
	defer conn.Close()

	cliente := pb.NewProductSaleServiceClient(conn)

	r := gin.Default()
	r.POST("/sale", func(ctx *gin.Context) {
		var venta SolicitudVenta
		if err := ctx.BindJSON(&venta); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Mapear string de Categoría a Enum de Protobuf
		catEnum, ok := pb.CategoriaProducto_value[venta.Categoria]
		if !ok {
			catEnum = int32(pb.CategoriaProducto_UNKNOWN)
		}

		// Preparar solicitud gRPC
		pbReq := &pb.ProductSaleRequest{
			Categoria:       pb.CategoriaProducto(catEnum),
			ProductoId:      venta.Producto,
			Precio:          venta.Precio,
			CantidadVendida: venta.CantidadVendida,
		}

		// Kafka puede tardar un poco al primer produce; damos más margen al gRPC.
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := cliente.ProcesarVenta(c, pbReq)
		if err != nil {
			log.Printf("Error llamando a gRPC: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "falló llamada gRPC"})
			return
		}

		log.Printf("Venta procesada: %+v, Respuesta: %s", venta, res.Estado)
		ctx.JSON(http.StatusOK, gin.H{"status": "procesado", "grpc_state": res.Estado})
	})

	log.Println("Iniciando API REST Go en :8080")
	r.Run(":8080")
}
