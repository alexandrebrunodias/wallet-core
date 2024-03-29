package gateway

import (
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/google/uuid"
)

type CustomerGateway interface {
	Create(customer *entity.Customer) error
	GetByID(ID uuid.UUID) (*entity.Customer, error)
}
