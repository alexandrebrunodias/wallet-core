package postgre

import (
	"database/sql"
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/google/uuid"
)

type CustomerPgGatewayDB struct {
	DB *sql.DB
}

func NewCustomerPgGateway(db *sql.DB) *CustomerPgGatewayDB {
	return &CustomerPgGatewayDB{DB: db}
}

func (c *CustomerPgGatewayDB) Save(customer *entity.Customer) error {
	stmt, err := c.DB.Prepare(
		"INSERT INTO customers (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(customer.ID, customer.Name, customer.Email, customer.CreatedAt, customer.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (c *CustomerPgGatewayDB) GetByID(ID uuid.UUID) (*entity.Customer, error) {
	customer := &entity.Customer{}
	query := `SELECT id, name, email, created_at, updated_at 
				FROM customers 
 				WHERE id = ?`

	stmt, err := c.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(ID).
		Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	return customer, nil
}
