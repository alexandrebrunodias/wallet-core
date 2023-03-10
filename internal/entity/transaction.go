package entity

import (
	"errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Status string

const (
	COMPLETED Status = "COMPLETED"
	FAILED           = "FAILED"
)

type Transaction struct {
	ID          uuid.UUID
	fromAccount *Account
	toAccount   *Account
	Status      Status
	Amount      decimal.Decimal
	CreatedAt   time.Time
}

func NewTransaction(fromAccount *Account, toAccount *Account, amount decimal.Decimal) (*Transaction, error) {
	transaction := &Transaction{
		ID:          uuid.New(),
		fromAccount: fromAccount,
		toAccount:   toAccount,
		Status:      FAILED,
		Amount:      amount,
		CreatedAt:   time.Now().UTC(),
	}
	err := transaction.Validate()
	if err != nil {
		return nil, err
	}
	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *Transaction) Commit() error {
	if t.Status == COMPLETED {
		return errors.New("transaction is already COMPLETED")
	}

	if err := t.fromAccount.Debit(t.Amount); err != nil {
		return err
	}
	if err := t.toAccount.Credit(t.Amount); err != nil {
		return err
	}

	t.Status = COMPLETED
	return nil
}

func (t *Transaction) Validate() error {
	if t.fromAccount == nil || t.toAccount == nil {
		return errors.New("neither 'fromAccount' nor 'toAccount' can be nil")
	}
	if t.Amount.IsNegative() || t.Amount.IsZero() {
		return errors.New("'amount' must be a non zero positive number")
	}
	return nil
}
