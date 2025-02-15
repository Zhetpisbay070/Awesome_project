package repository

import (
	"awesomeProject1/internal/entity"

	"context"
	"database/sql"
	"time"
)

type postgresRepository struct {
	db *sql.DB
}

func (r *postgresRepository) CreateOrder(ctx context.Context, order *entity.Order) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO orders(id, user_id, created_at, updated_at, delivery_deadline, price, delivery_type, address, order_status) 
  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		order.ID,
		order.UserID,
		order.CreatedAt.Format(time.RFC3339),
		order.UpdatedAt.Format(time.RFC3339),
		order.DeliveryDeadLine.Format(time.RFC3339),
		order.Price,
		order.DeliveryType,
		order.Address,
		order.OrderStatus)

	if err != nil {
		return err
	}

	return nil
}
func (r *postgresRepository) GetOrderByID(ctx context.Context, id string) (*entity.Order, error) {
	var order entity.Order

	row := r.db.QueryRowContext(ctx, `SELECT id, user_id, price FROM orders WHERE id = ?`, id)

	err := row.Scan(&order.ID, &order.UserID, &order.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}

func (r *postgresRepository) ProductExist(ctx context.Context, productID string) (bool, error) {
	var exist bool
	row := r.db.QueryRowContext(ctx, `select exists(select 1 from products WHERE id = ?)`, productID)

	err := row.Scan(&exist)
	if err != nil {
		return false, nil
	}

	return exist, nil
}

func (r *postgresRepository) UpdateOrder(ctx context.Context, order *entity.Order) error {
	_, err := r.db.ExecContext(ctx, `UPDATE orders SET user_id = ?, updated_at = ?, delivery_deadline = ?, price = ?, delivery_type = ?, address = ?, order_status = ?
WHERE id = ?`,
		order.UserID,
		order.UpdatedAt.Format(time.RFC3339),
		order.DeliveryDeadLine.Format(time.RFC3339),
		order.Price,
		order.DeliveryType,
		order.Address,
		order.OrderStatus,
		order.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
func (r *postgresRepository) GetOrders(ctx context.Context, req *entity.GetOrders) ([]entity.Order, error) {
	var orders []entity.Order

	_, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, price, delivery_deadline, delivery_type, address, order_status
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC`, req.UserID)

	if err != nil {
		return nil, err

	}

	return orders, nil
}
