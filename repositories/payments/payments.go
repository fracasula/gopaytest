package payments

import (
	"fmt"
	"gopaytest/models/payment"
	"sort"
	"sync"
)

// UpdateError is used as a return value by the UpdatePayment method to detect whether the error is
// about a version conflict or not
//go:generate counterfeiter . UpdateError
type UpdateError interface {
	Error() string
	IsConflict() bool
}

// Repository represents a contract for querying payments from an arbitrary data source
//go:generate counterfeiter . Repository
type Repository interface {
	// CreatePayment stores the payment in the database and returns error on connection errors
	CreatePayment(p payment.Payment) error
	// CountPayments returns the total number of payments in the database or error on connection errors
	CountPayments() (int, error)
	// GetPayments returns a slice of payments or error on connection errors
	GetPayments(limit, offset int) ([]payment.Payment, error)
	// GetPayment returns a *Payment if found or nil if not found and error on connection errors
	GetPayment(ID string) (*payment.Payment, error)
	// UpdatePayment returns true if found or false if not found and error on connection errors
	UpdatePayment(p payment.Payment) (bool, UpdateError)
	// DeletePayment returns true|false respectively if the record was found and deleted or not found and error on
	// connection errors
	DeletePayment(ID string) (bool, error)
}

type updateError struct {
	message       string
	conflictError bool
}

func (e *updateError) Error() string {
	return fmt.Sprintf("Update error: %v (conflict %v)", e.message, e.conflictError)
}

func (e *updateError) IsConflict() bool { return e.conflictError }

type paymentID = string
type inMemoryPaymentsRepo struct {
	lock       *sync.RWMutex
	store      map[paymentID]payment.Payment
	sortedKeys []paymentID // we need a sorted set to avoid breaking pagination
}

func (r *inMemoryPaymentsRepo) CreatePayment(p payment.Payment) error {
	defer r.lock.Unlock()
	r.lock.Lock()

	_, ok := r.store[p.ID]
	if ok {
		// value already exists, overwrite it without updating the keys
		r.store[p.ID] = p
		return nil
	}

	// value doesn't exist, store it, add the key and sort the keys slice
	r.store[p.ID] = p
	r.sortedKeys = append(r.sortedKeys, p.ID)
	sort.Strings(r.sortedKeys)

	return nil
}

func (r *inMemoryPaymentsRepo) CountPayments() (int, error) {
	defer r.lock.RUnlock()
	r.lock.RLock()

	// error is nil here since it's an in-memory repository, there can't be any connection error
	return len(r.store), nil
}

func (r *inMemoryPaymentsRepo) GetPayments(limit, offset int) ([]payment.Payment, error) {
	defer r.lock.RUnlock()
	r.lock.RLock()

	count := 0
	values := make([]payment.Payment, 0)

	for i, id := range r.sortedKeys {
		if i >= offset && count < limit {
			values = append(values, r.store[id])
			count++
		}
	}

	// error is nil here since it's an in-memory repository, there can't be any connection error
	return values, nil
}

func (r *inMemoryPaymentsRepo) GetPayment(ID string) (*payment.Payment, error) {
	defer r.lock.RUnlock()
	r.lock.RLock()

	p, ok := r.store[ID]
	if ok {
		return &p, nil
	}

	// error is nil here since it's an in-memory repository, there can't be any connection error
	return nil, nil
}

func (r *inMemoryPaymentsRepo) UpdatePayment(p payment.Payment) (bool, UpdateError) {
	defer r.lock.Unlock()
	r.lock.Lock()

	existingPayment, ok := r.store[p.ID]
	if !ok {
		return false, nil // not found
	}

	if existingPayment.ID != p.ID {
		return false, &updateError{
			message: fmt.Sprintf(
				"it is not possible to change the ID of an existing payment (was %q, got %q instead)",
				existingPayment.ID, p.ID,
			),
			conflictError: false,
		}
	}

	if p.Version != existingPayment.Version {
		return false, &updateError{
			message: fmt.Sprintf(
				"could not update payment because of version conflict (got %d, expected %d)",
				p.Version, existingPayment.Version,
			),
			conflictError: true,
		}
	}

	p.Version = existingPayment.Version + 1
	r.store[p.ID] = p

	// error is nil here since it's an in-memory repository, there can't be any connection error
	return true, nil
}

func (r *inMemoryPaymentsRepo) DeletePayment(ID string) (bool, error) {
	defer r.lock.Unlock()
	r.lock.Lock()

	_, ok := r.store[ID]
	if ok {
		newSortedKeys := make([]paymentID, len(r.sortedKeys)-1)

		index := 0
		for _, key := range r.sortedKeys {
			if key == ID {
				continue
			}
			newSortedKeys[index] = key
			index++
		}

		delete(r.store, ID)
		r.sortedKeys = newSortedKeys
	}

	return ok, nil
}

// NewInMemoryPaymentsRepository returns an in-memory implementation of the Repository interface
func NewInMemoryPaymentsRepository() Repository {
	return &inMemoryPaymentsRepo{
		lock:       &sync.RWMutex{},
		store:      make(map[paymentID]payment.Payment),
		sortedKeys: make([]paymentID, 0),
	}
}
