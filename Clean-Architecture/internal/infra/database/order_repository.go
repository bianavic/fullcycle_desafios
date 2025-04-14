package database

import (
	"database/sql"
	"fmt"
	"github.com/bianavic/fullcycle_clean-architecture/internal/entity"
	"time"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	smt, err := r.DB.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = smt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *OrderRepository) List() ([]entity.Order, error) {
	rows, err := r.DB.Query("SELECT id, price, tax, final_price, created_at FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		var createdAtRaw interface{}
		err := rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice, &createdAtRaw)
		if err != nil {
			return nil, err
		}

		switch v := createdAtRaw.(type) {
		case time.Time:
			order.CreatedAt = v
		case []byte:
			order.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(v))
			if err != nil {
				return nil, fmt.Errorf("failed to parse created_at: %w", err)
			}
		default:
			return nil, fmt.Errorf("unsupported type for created_at: %T", v)
		}

		orders = append(orders, order)
	}
	return orders, nil
}
