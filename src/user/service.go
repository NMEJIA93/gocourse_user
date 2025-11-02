package user

import (
	"github.com/NMEJIA93/gocourse_domain/domain"
	"log"
)

type (
	Service interface {
		Create(dto CreateUserDTO) (*domain.User, error)
		Get(id string) (*domain.User, error)
		GetAll(filter Filters, offset int, limit int) ([]domain.User, error)
		Delete(id string) error
		Update(id string, firstName *string, lastName *string, email *string, phone *string) error
		Count(filter Filters) (int, error)
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

func (s service) Create(dto CreateUserDTO) (*domain.User, error) {
	s.log.Println("Create User Service")
	user := domain.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Phone:     dto.Phone,
	}
	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s service) GetAll(filter Filters, offset int, limit int) ([]domain.User, error) {
	users, err := s.repo.GetAll(filter, offset, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s service) Get(id string) (*domain.User, error) {
	user, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s service) Update(id string, firstName *string, lastName *string, email *string, phone *string) error {
	return s.repo.Update(id, firstName, lastName, email, phone)
}

func (s service) Count(filter Filters) (int, error) {
	return s.repo.Count(filter)
}
