package create_transaction

import (
	"context"
	"errors"
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/alexandrebrunodias/wallet-core/pkg/events"
	"github.com/alexandrebrunodias/wallet-core/pkg/uow"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	m "github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateTransactionUseCase_Execute_CreateSuccessfully(t *testing.T) {
	fromCustomer, _ := entity.NewCustomer("fromCustomer", "alexandrebrunodias@gmail.com")
	expectedFromAccount, _ := entity.NewAccount(fromCustomer)
	_ = expectedFromAccount.Credit(decimal.NewFromInt(2000))

	toCustomer, _ := entity.NewCustomer("toCustomer", "alexandrebrunodias@gmail.com")
	expectedToAccount, _ := entity.NewAccount(toCustomer)

	expectedAmount := decimal.NewFromInt(1000)

	command := CreateTransactionCommand{
		FromAccountID: expectedFromAccount.ID,
		ToAccountID:   expectedToAccount.ID,
		Amount:        expectedAmount,
	}

	unitOfWorkMock := &UnitOfWorkMock{}
	unitOfWorkMock.On("Do", m.Anything, m.Anything).Return(nil)

	eventPublisherMock := &EventPublisherMock{}
	eventPublisherMock.On("Register", m.Anything).Return(eventPublisherMock)
	eventPublisherMock.On("Publish", m.Anything)

	useCase := NewCreateTransactionUseCase(unitOfWorkMock, eventPublisherMock)
	output, err := useCase.Execute(context.Background(), command)

	assert.Nil(t, err)
	assert.NotNil(t, output.ID)

	unitOfWorkMock.AssertExpectations(t)
	unitOfWorkMock.AssertNumberOfCalls(t, "Do", 1)

	eventPublisherMock.AssertExpectations(t)
	eventPublisherMock.AssertNumberOfCalls(t, "Register", 1)
	eventPublisherMock.AssertNumberOfCalls(t, "Publish", 1)
}

func TestCreateTransactionUseCase_Execute_FailDueToErrorOnUnitOfWorkTransaction(t *testing.T) {
	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	expectedErrorMessage := "any error on unit of work"

	command := CreateTransactionCommand{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        decimal.NewFromInt(1000),
	}

	unitOfWorkMock := &UnitOfWorkMock{}
	unitOfWorkMock.
		On("Do", m.Anything, m.Anything).
		Return(errors.New(expectedErrorMessage))

	eventPublisherMock := &EventPublisherMock{}

	useCase := NewCreateTransactionUseCase(unitOfWorkMock, eventPublisherMock)
	output, err := useCase.Execute(context.Background(), command)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedErrorMessage, err.Error())

	unitOfWorkMock.AssertExpectations(t)
	unitOfWorkMock.AssertNumberOfCalls(t, "Do", 1)

	eventPublisherMock.AssertExpectations(t)
	eventPublisherMock.AssertNotCalled(t, "Register")
	eventPublisherMock.AssertNotCalled(t, "Publish")
}

type EventPublisherMock struct {
	m.Mock
}

func (e *EventPublisherMock) Register(event events.Event) events.EventPublisherInterface {
	args := e.Called(event)
	return args.Get(0).(events.EventPublisherInterface)
}

func (e *EventPublisherMock) Publish() {
	e.Called()
}

type UnitOfWorkMock struct {
	m.Mock
}

func (m *UnitOfWorkMock) Do(_ context.Context, fn func(unitOfWork *uow.UnitOfWork) error) error {
	args := m.Called(fn)
	return args.Error(0)
}

func (m *UnitOfWorkMock) Add(name string, repository uow.Repository) {}

func (m *UnitOfWorkMock) Remove(name string) {}

func (m *UnitOfWorkMock) GetRepository(ctx context.Context, name string) (interface{}, error) {
	return nil, nil
}

func (m *UnitOfWorkMock) CommitOrRollback() error {
	return nil
}

func (m *UnitOfWorkMock) RollBack() error {
	return nil
}
