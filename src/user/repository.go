package user

import (
	"context"
	"fmt"
	"github.com/NMEJIA93/gocourse_domain/domain"
	"log"
	"strings"

	"gorm.io/gorm"
)

type Respository interface {
	Create(ctx context.Context, user *domain.User) error
	GetAll(ctx context.Context, filters Filters, offset int, limit int) ([]domain.User, error)
	Get(ctx context.Context, id string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, firstName *string, lasName *string, email *string, phone *string) error
	Count(ctx context.Context, filters Filters) (int, error)
}

type repository struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepository(log *log.Logger, db *gorm.DB) Respository {
	return &repository{
		log: log,
		db:  db,
	}
}

func (r *repository) Create(ctx context.Context, user *domain.User) error {

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		r.log.Printf("Error while creating user: %v", err)
		return err
	}

	r.log.Println("user Created with id: ", user.ID)
	return nil
}

func (r *repository) GetAll(ctx context.Context, filter Filters, offset int, limit int) ([]domain.User, error) {
	var user []domain.User

	tx := r.db.WithContext(ctx).Model(&domain.User{})
	tx = applyFilters(tx, filter)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("Created_at desc").Find(&user)
	if result.Error != nil {
		r.log.Println(result.Error)
		return nil, result.Error
	}
	return user, nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.User, error) {
	user := domain.User{ID: id}

	err := r.db.WithContext(ctx).First(&user).Error

	if err != nil {
		r.log.Println(err)
		return nil, err
	}
	return &user, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	user := domain.User{ID: id}
	//Eliminado Fisico
	//result := r.db.Delete(&user)
	err := r.db.WithContext(ctx).First(&user).Error
	if err != nil {
		return err
	}

	result := r.db.Delete(&user).Error

	if result != nil {
		r.log.Println(result)
		return result
	}

	return nil
}

func (r *repository) Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) error {
	values := make(map[string]interface{})

	if firstName != nil {
		values["first_name"] = *firstName
	}
	if lastName != nil {
		values["last_name"] = *lastName
	}
	if email != nil {
		values["email"] = *email
	}
	if phone != nil {
		values["phone"] = *phone
	}

	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(values).Error; err != nil {
		r.log.Println(err)
		return err
	}

	return nil
}

func (r *repository) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := r.db.WithContext(ctx).Model(domain.User{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		r.log.Println(err)
		return 0, err
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.FirstName != "" {
		filters.FirstName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.FirstName))
		tx = tx.Where("lower(first_name) LIKE ?", filters.FirstName)
	}

	if filters.LastName != "" {
		filters.LastName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.LastName))
		tx = tx.Where("lower(first_name) LIKE ?", filters.LastName)
	}
	return tx
}
