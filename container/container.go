package container

import (
	"gopaytest/repositories/payments"
	"log"
	"os"
)

//go:generate counterfeiter . Container
type Container interface {
	BaseURL() string
	PaymentsRepository() payments.Repository
	Logger() *log.Logger
}

type container struct {
	baseURL            string
	paymentsRepository payments.Repository
	logger             *log.Logger
}

func (c *container) BaseURL() string {
	return c.baseURL
}

func (c *container) PaymentsRepository() payments.Repository {
	return c.paymentsRepository
}

func (c *container) Logger() *log.Logger {
	return c.logger
}

func NewContainer(baseURL string) Container {
	return &container{
		baseURL:            baseURL,
		paymentsRepository: payments.NewInMemoryPaymentsRepository(),
		logger:             log.New(os.Stdout, "", log.LstdFlags),
	}
}
