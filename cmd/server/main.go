// cmd/server/main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"

	g "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc"
	pb "github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/grpc/proto"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/redis"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/repository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("failed to load env variables", err)
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	log.Println(redisAddr)
	redisUsername := os.Getenv("REDIS_USERNAME")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0 // Default DB
	cache := redis.NewCache(redisAddr, redisUsername, redisPassword, redisDB, 5*time.Minute)
	log.Println(cache)
	if err := cache.Ping(context.Background()); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	initDB(db)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("321dsaf"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash default user password: %v", err)
	} else {
		_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2) ON CONFLICT (username) DO NOTHING",
			"01901901901@mailinator.com", string(hashedPassword))
		if err != nil {
			log.Printf("failed to insert default user: %v", err)
		}
	}

	repo := repository.NewPostgresRepository(db)
	srv := g.NewServer(repo, cache)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(g.AuthInterceptor))
	pb.RegisterOrderServiceServer(grpcServer, srv)

	fmt.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initDB(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			consignment_id VARCHAR(255) PRIMARY KEY,
			created_at TIMESTAMP NOT NULL,
			description TEXT,
			merchant_order_id VARCHAR(255),
			recipient_name VARCHAR(255) NOT NULL,
			recipient_address TEXT NOT NULL,
			recipient_phone VARCHAR(20) NOT NULL,
			order_amount FLOAT NOT NULL,
			total_fee FLOAT NOT NULL,
			instruction TEXT,
			order_type_id BIGINT NOT NULL,
			cod_fee FLOAT NOT NULL,
			promo_discount FLOAT NOT NULL,
			discount FLOAT NOT NULL,
			delivery_fee FLOAT NOT NULL,
			status VARCHAR(50) NOT NULL,
			order_type VARCHAR(50) NOT NULL,
			item_type BIGINT NOT NULL,
			store_name VARCHAR(255),
			store_contact_phone VARCHAR(20),
			cod_amount FLOAT NOT NULL,
			delivery_charge FLOAT NOT NULL,
			user_id BIGINT REFERENCES users(id),
			store_id BIGINT NOT NULL,
			recipient_city BIGINT NOT NULL,
			recipient_zone BIGINT NOT NULL,
			recipient_area BIGINT NOT NULL,
			delivery_type BIGINT NOT NULL,
			item_quantity BIGINT NOT NULL,
			item_weight FLOAT NOT NULL,
			amount_to_collect FLOAT NOT NULL
		)`,
	}
	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			log.Fatalf("failed to init DB: %v", err)
		}
	}
}
