package infrastructure

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func InitDB(db *sql.DB) {
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

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("321dsaf"), bcrypt.DefaultCost)
	_, _ = db.Exec(`INSERT INTO users (username, password) VALUES ($1, $2)
	                ON CONFLICT (username) DO NOTHING`,
		"01901901901@mailinator.com", string(hashedPassword))
}
