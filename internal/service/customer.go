package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yusriltakeuchi/golang-perpustakaan/domain"
	"github.com/yusriltakeuchi/golang-perpustakaan/dto"
)

type customerService struct {
	customerRepository domain.CustomerRepository
}

func NewCustomer(customerRepository domain.CustomerRepository) domain.CustomerService {
	return &customerService{
		customerRepository: customerRepository,
	}
}

// Index implements [domain.CustomerService].
func (c *customerService) Index(ctx context.Context) ([]dto.CustomerData, error) {
	customers, err := c.customerRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var customerData []dto.CustomerData
	for _, v := range customers {
		customerData = append(customerData, dto.CustomerData{
			ID:   v.ID,
			Name: v.Name,
			Code: v.Code,
		})
	}
	return customerData, nil
}

// Create implements [domain.CustomerService].
func (c *customerService) Create(ctx context.Context, req dto.CreateCustomerRequest) error {
	customer := domain.Customer{
		ID:        uuid.NewString(),
		Code:      req.Code,
		Name:      req.Name,
		CreatedAt: sql.NullTime{Valid: true, Time: time.Now()},
	}
	return c.customerRepository.Save(ctx, &customer)
}

// Update implements [domain.CustomerService].
func (c *customerService) Update(ctx context.Context, req dto.UpdateCustomerRequest) error {
	persisted, err := c.customerRepository.FindById(ctx, req.ID)
	if err != nil {
		return err
	}
	if persisted.ID == "" {
		return errors.New("data customer tidak ditemukan")
	}
	persisted.Name = req.Name
	persisted.Code = req.Code
	persisted.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	return c.customerRepository.Update(ctx, &persisted)
}

// Delete implements [domain.CustomerService].
func (c *customerService) Delete(ctx context.Context, id string) error {
	exists, err := c.customerRepository.FindById(ctx, id)
	if err != nil {
		return err
	}
	if exists.ID == "" {
		return errors.New("data customer tidak ditemukan")
	}
	return c.customerRepository.Delete(ctx, id)
}

// Show implements [domain.CustomerService].
func (c *customerService) Show(ctx context.Context, id string) (dto.CustomerData, error) {
	persisted, err := c.customerRepository.FindById(ctx, id)
	if err != nil {
		return dto.CustomerData{}, err
	}
	if persisted.ID == "" {
		return dto.CustomerData{}, errors.New("data customer tidak ditemukan")
	}
	return dto.CustomerData{
		ID:   persisted.ID,
		Code: persisted.Code,
		Name: persisted.Name,
	}, nil
}
