package repository

import (
	"UserService/src/model"
)

type RoleRepository interface {
	Create(role *model.Role) (*model.Role, error)
	FindById(id string) (*model.Role, error)
	FindByCode(code string) (*model.Role, error)
	FindAll() ([]*model.Role, error)
	Update(role *model.Role) (*model.Role, error)
	Delete(id string) error
	FindAllByPageAndSize(page, size int) ([]*model.Role, error)
	Count() (int64, error)
	FindByIds(ids []string) ([]model.Role, error)
}