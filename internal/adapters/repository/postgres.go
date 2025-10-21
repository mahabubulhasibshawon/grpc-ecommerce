// internal/adapters/repository/postgres.go
package repository

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/domain"
	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/ports"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) ports.OrderRepositoryPort {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) FindUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (
			consignment_id, created_at, description, merchant_order_id, recipient_name, recipient_address, recipient_phone,
			order_amount, total_fee, instruction, order_type_id, cod_fee, promo_discount, discount, delivery_fee, status,
			order_type, item_type, store_name, store_contact_phone, cod_amount, delivery_charge, user_id, store_id,
			recipient_city, recipient_zone, recipient_area, delivery_type, item_quantity, item_weight, amount_to_collect
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31)
	`
	_, err := r.db.ExecContext(ctx, query,
		order.ConsignmentID, order.CreatedAt, order.Description, order.MerchantOrderID, order.RecipientName, order.RecipientAddress, order.RecipientPhone,
		order.OrderAmount, order.TotalFee, order.Instruction, order.OrderTypeID, order.CODFee, order.PromoDiscount, order.Discount, order.DeliveryFee, order.Status,
		order.OrderType, order.ItemType, order.StoreName, order.StoreContactPhone, order.CODAmount, order.DeliveryCharge, order.UserID, order.StoreID,
		order.RecipientCity, order.RecipientZone, order.RecipientArea, order.DeliveryType, order.ItemQuantity, order.ItemWeight, order.AmountToCollect,
	)
	return err
}

func (r *PostgresRepository) ListOrders(ctx context.Context, userID int64, limit, page int64) ([]*domain.Order, int64, error) {
	if limit < 1 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var total int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM orders WHERE user_id = $1", userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	query := `
		SELECT consignment_id, created_at, description, merchant_order_id, recipient_name, recipient_address, recipient_phone,
			order_amount, total_fee, instruction, order_type_id, cod_fee, promo_discount, discount, delivery_fee, status,
			order_type, item_type, store_name, store_contact_phone, cod_amount, delivery_charge, store_id,
			recipient_city, recipient_zone, recipient_area, delivery_type, item_quantity, item_weight, amount_to_collect
		FROM orders WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		o := &domain.Order{UserID: userID}
		err := rows.Scan(
			&o.ConsignmentID, &o.CreatedAt, &o.Description, &o.MerchantOrderID, &o.RecipientName, &o.RecipientAddress, &o.RecipientPhone,
			&o.OrderAmount, &o.TotalFee, &o.Instruction, &o.OrderTypeID, &o.CODFee, &o.PromoDiscount, &o.Discount, &o.DeliveryFee, &o.Status,
			&o.OrderType, &o.ItemType, &o.StoreName, &o.StoreContactPhone, &o.CODAmount, &o.DeliveryCharge, &o.StoreID,
			&o.RecipientCity, &o.RecipientZone, &o.RecipientArea, &o.DeliveryType, &o.ItemQuantity, &o.ItemWeight, &o.AmountToCollect,
		)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}
	return orders, total, nil
}

func (r *PostgresRepository) CancelOrder(ctx context.Context, consignmentID string, userID int64) error {
	res, err := r.db.ExecContext(ctx, "UPDATE orders SET status = 'Cancelled' WHERE consignment_id = $1 AND user_id = $2 AND status = 'Pending'", consignmentID, userID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("order not found, unauthorized, or cannot cancel")
	}
	return nil
}
