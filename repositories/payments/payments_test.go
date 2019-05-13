package payments

import (
	"gopaytest/models/payment"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInMemoryPaymentsRepo_CountPayments(t *testing.T) {
	repo := getMockedRepo()
	require.Nil(t, repo.CreatePayment(payment.MockedPayment("payment 001")))
	require.Nil(t, repo.CreatePayment(payment.MockedPayment("payment 002")))
	require.Nil(t, repo.CreatePayment(payment.MockedPayment("payment 003")))
	require.Nil(t, repo.CreatePayment(payment.MockedPayment("payment 004")))

	count, err := repo.CountPayments()
	require.Nil(t, err)
	require.Equal(t, 4, count)
}

func TestInMemoryPaymentsRepo_GetPayments(t *testing.T) {
	repo := getMockedRepo()
	require.Nil(t, repo.CreatePayment(payment.MockedPayment("payment 001")))
	require.Nil(t, repo.CreatePayment(payment.MockedPayment("payment 002")))
	require.Nil(t, repo.CreatePayment(payment.MockedPayment("payment 003")))
	require.Nil(t, repo.CreatePayment(payment.MockedPayment("payment 004")))

	page1, err := repo.GetPayments(2, 0)
	require.Nil(t, err)
	require.Equal(t, []payment.Payment{
		repo.store["payment 001"],
		repo.store["payment 002"],
	}, page1)

	page2, err := repo.GetPayments(10, 2)
	require.Nil(t, err)
	require.Equal(t, []payment.Payment{
		repo.store["payment 003"],
		repo.store["payment 004"],
	}, page2)

	page3, err := repo.GetPayments(10, 10)
	require.Nil(t, err)
	require.Equal(t, []payment.Payment{}, page3)
}

func TestInMemoryPaymentsRepo_GetPayment(t *testing.T) {
	repo := getMockedRepo()
	require.Nil(t, repo.CreatePayment(payment.MockedPayment("payment 001")))
	require.Nil(t, repo.CreatePayment(payment.MockedPayment("payment 002")))
	require.Nil(t, repo.CreatePayment(payment.MockedPayment("payment 003")))
	require.Nil(t, repo.CreatePayment(payment.MockedPayment("payment 004")))

	p1, err := repo.GetPayment("payment 001")
	require.Nil(t, err)
	require.Equal(t, repo.store["payment 001"], *p1)

	p3, err := repo.GetPayment("payment 003")
	require.Nil(t, err)
	require.Equal(t, repo.store["payment 003"], *p3)
}

func TestInMemoryPaymentsRepo_UpdatePayment(t *testing.T) {
	repo := getMockedRepo()
	p1 := payment.MockedPayment("payment 001")
	p2 := payment.MockedPayment("payment 002")
	p3 := payment.MockedPayment("payment 003")
	require.Nil(t, repo.CreatePayment(p1))
	require.Nil(t, repo.CreatePayment(p2))
	require.Nil(t, repo.CreatePayment(p3))

	require.Equal(t, 1, p2.Version, "Version expected to be 1 before updating")
	p2.OrganisationID = "a new organisation id"
	ok, updateErr := repo.UpdatePayment(p2)
	require.Nil(t, updateErr)
	require.True(t, ok)
	require.Equal(t, 1, p2.Version, "Version on original objected expected to be 1 after updating")

	updatedP2, err := repo.GetPayment(p2.ID)
	require.Nil(t, err)
	require.Equal(t, p2.Version+1, updatedP2.Version, "Version expected to be bumped by 1 after update")

	nonExistentP := payment.MockedPayment("non-existent")
	ok, err = repo.UpdatePayment(nonExistentP)
	require.Nil(t, err)
	require.False(t, ok)

	// try to update again using the old payment, we should get a conflict
	ok, updateErr = repo.UpdatePayment(p2)
	require.False(t, ok)
	require.Equal(t, &updateError{
		message:       "could not update payment because of version conflict (got 1, expected 2)",
		conflictError: true,
	}, updateErr)
}

func TestInMemoryPaymentsRepo_DeletePayment(t *testing.T) {
	repo := getMockedRepo()
	require.Nil(t, repo.CreatePayment(payment.MockedPayment("payment 001")))
	require.Nil(t, repo.CreatePayment(payment.MockedPayment("payment 002")))
	require.Nil(t, repo.CreatePayment(payment.MockedPayment("payment 003")))
	require.Nil(t, repo.CreatePayment(payment.MockedPayment("payment 004")))
	require.Nil(t, repo.CreatePayment(payment.MockedPayment("payment 005")))

	ok, err := repo.DeletePayment("payment 002")
	require.Nil(t, err)
	require.True(t, ok)

	ok, err = repo.DeletePayment("payment 004")
	require.Nil(t, err)
	require.True(t, ok)

	count, err := repo.CountPayments()
	require.Nil(t, err)
	require.Equal(t, 3, count)

	page1, err := repo.GetPayments(10, 0)
	require.Nil(t, err)
	require.Equal(t, []payment.Payment{
		repo.store["payment 001"],
		repo.store["payment 003"],
		repo.store["payment 005"],
	}, page1)

	ok, err = repo.DeletePayment("non-existent")
	require.Nil(t, err)
	require.False(t, ok)
}

func getMockedRepo() *inMemoryPaymentsRepo {
	return &inMemoryPaymentsRepo{
		lock:  &sync.RWMutex{},
		store: make(map[paymentID]payment.Payment),
	}
}
