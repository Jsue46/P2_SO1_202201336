package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type SolicitudVenta struct {
	Categoria       string  `json:"categoria"`
	Producto        string  `json:"producto"`
	Precio          float64 `json:"precio"`
	CantidadVendida int32   `json:"cantidad_vendida"`
}

func main() {
	addrKafka := os.Getenv("KAFKA_BROKER")
	if addrKafka == "" {
		addrKafka = "my-cluster-kafka-bootstrap:9092"
	}

	addrRedis := os.Getenv("VALKEY_ADDR")
	if addrRedis == "" {
		addrRedis = "valkey-service:6379"
	}

	// Conectar a Redis (Valkey)
	clienteRedis := redis.NewClient(&redis.Options{
		Addr: addrRedis,
	})

	// Conectar a Kafka
	lector := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{addrKafka},
		Topic:    "sales-topic",
		GroupID:  "sales-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer lector.Close()

	log.Println("Consumidor iniciado...")

	ctx := context.Background()

	for {
		m, err := lector.ReadMessage(ctx)
		if err != nil {
			log.Printf("no se pudo leer mensaje: %v", err)
			break
		}

		// Decodificar mensaje
		var req SolicitudVenta
		if err := json.Unmarshal(m.Value, &req); err != nil {
			log.Printf("Error decodificando: %v", err)
			continue
		}

		// --- Lógica de Redis (Valkey) ---

		// 1. Contador por Categoría (BARRAS: Cantidad total de reportes por categoría)
		nuevoConteo, errConteo := clienteRedis.HIncrBy(ctx, "stats:conteo_categoria", req.Categoria, 1).Result()
		if errConteo != nil {
			log.Printf("Error actualizando conteo: %v", errConteo)
		}

		// 2. Suma de Precios por Categoría (Para sacar Promedios)
		nuevaSuma, errSuma := clienteRedis.HIncrByFloat(ctx, "stats:suma_precio_categoria", req.Categoria, req.Precio).Result()
		if errSuma != nil {
			log.Printf("Error actualizando suma: %v", errSuma)
		}

		// 2.1. Suma de Cantidad de Productos (Para Promedio de Productos)
		nuevaSumaCantidad, errCant := clienteRedis.HIncrBy(ctx, "stats:suma_cantidad_categoria", req.Categoria, int64(req.CantidadVendida)).Result()
		if errCant != nil {
			log.Printf("Error actualizando suma cantidad: %v", errCant)
		}

		// 2.5 CÁLCULO DE PROMEDIOS (Backend)
		if errConteo == nil && nuevoConteo > 0 {
			// A. Precio Promedio
			if errSuma == nil {
				promedioPrecio := nuevaSuma / float64(nuevoConteo)
				clienteRedis.HSet(ctx, "stats:promedio_precio_categoria", req.Categoria, promedioPrecio)
			}

			// B. Producto Promedio (Cantidad Promedio por Transacción)
			if errCant == nil {
				promedioProds := float64(nuevaSumaCantidad) / float64(nuevoConteo)
				clienteRedis.HSet(ctx, "stats:promedio_productos_categoria", req.Categoria, promedioProds)
			}
		}

		// 3. Max/Min Precio Global (TARJETAS: Precio más alto/bajo registrado)
		registrarMaxMinGlobal(ctx, clienteRedis, req.Precio)

		// 4. Rankings de Productos (Para Top/Menos vendidos)
		// Global:
		clienteRedis.ZIncrBy(ctx, "stats:ranking_productos_global", 1, req.Producto)
		// Por Categoría (Para la Sección Carnet):
		if req.Categoria == "Hogar" {
			clienteRedis.ZIncrBy(ctx, "stats:ranking_productos_hogar", 1, req.Producto)

			// 5. Time Series (Sección Carnet: Variación precio producto Hogar)
			// Usamos un Hash temporal para guardar el último precio del producto específico vendido
			clienteRedis.XAdd(ctx, &redis.XAddArgs{
				Stream: "stats:stream:Hogar",
				Values: map[string]interface{}{
					"producto": req.Producto,
					"precio":   fmt.Sprintf("%.2f", req.Precio),
				},
			})
		}

		log.Printf("Venta procesada: %s - %s ($%.2f)", req.Categoria, req.Producto, req.Precio)
	}
}

func registrarMaxMinGlobal(ctx context.Context, rdb *redis.Client, precio float64) {
	// Revisar Max
	maxActualStr, _ := rdb.Get(ctx, "stats:max_precio").Result()
	maxActual, _ := strconv.ParseFloat(maxActualStr, 64)
	if precio > maxActual {
		rdb.Set(ctx, "stats:max_precio", precio, 0)
	}

	// Revisar Min
	minActualStr, err := rdb.Get(ctx, "stats:min_precio").Result()
	if err == redis.Nil {
		// Primer valor
		rdb.Set(ctx, "stats:min_precio", precio, 0)
	} else {
		minActual, _ := strconv.ParseFloat(minActualStr, 64)
		if precio < minActual {
			rdb.Set(ctx, "stats:min_precio", precio, 0)
		}
	}
}
