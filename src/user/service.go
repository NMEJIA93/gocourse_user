package user

import (
	"context"
	"github.com/NMEJIA93/gocourse_domain/domain"
	"log"
)

type (
	Service interface {
		Create(ctx context.Context, dto CreateUserDTO) (*domain.User, error)
		Get(ctx context.Context, id string) (*domain.User, error)
		GetAll(ctx context.Context, filter Filters, offset int, limit int) ([]domain.User, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) error
		Count(ctx context.Context, filter Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Respository
	}

	Filters struct {
		FirstName string
		LastName  string
	}
)

func NewService(log *log.Logger, repo Respository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, dto CreateUserDTO) (*domain.User, error) {
	s.log.Println("Create User Service")
	user := domain.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Phone:     dto.Phone,
	}
	if err := s.repo.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s service) GetAll(ctx context.Context, filter Filters, offset int, limit int) ([]domain.User, error) {
	users, err := s.repo.GetAll(ctx, filter, offset, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s service) Get(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s service) Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) error {
	return s.repo.Update(ctx, id, firstName, lastName, email, phone)
}

func (s service) Count(ctx context.Context, filter Filters) (int, error) {
	return s.repo.Count(ctx, filter)
}
