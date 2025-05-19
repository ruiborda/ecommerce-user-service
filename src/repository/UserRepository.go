package repository

import (
	"github.com/ruiborda/ecommerce-user-service/src/model"
)

type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	FindById(id string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindAll() ([]*model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(id string) error
	FindAllByPageAndSize(page, size int) ([]*model.User, error)
	Count() (int64, error)
	FindByIds(ids []string) ([]*model.User, error)
}
