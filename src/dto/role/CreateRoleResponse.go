package role

import "github.com/ruiborda/ecommerce-user-service/src/model"

type CreateRoleResponse struct {
	Id          string              `json:"id"`
	Code        string              `json:"code"`
	Permissions *[]model.Permission `json:"permissions"`
}
