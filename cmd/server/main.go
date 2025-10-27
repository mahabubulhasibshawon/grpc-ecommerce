// cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	g "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc"
	pb "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc/proto"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/repository"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/config"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/infrastructure"
)

func main() {
	cfg := config.Load()

	db := infrastructure.ConnectPostgres(cfg.DSN())
	defer db.Close()

	cache := infrastructure.ConnectRedis(cfg.RedisAddr, cfg.RedisUsername, cfg.RedisPassword, cfg.RedisTTL)

	infrastructure.InitDB(db)

	repo := repository.NewPostgresRepository(db)
	srv := g.NewServer(repo, cache)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(g.AuthInterceptor))
	pb.RegisterOrderServiceServer(grpcServer, srv)

	log.Printf("🚀 gRPC server listening on :%s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
