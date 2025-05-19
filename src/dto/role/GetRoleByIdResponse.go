package role

import "github.com/ruiborda/ecommerce-user-service/src/model"

type GetRoleByIdResponse struct {
	Id          string              `json:"id"`
	Code        string              `json:"code"`
	Permissions *[]model.Permission `json:"permissions"`
}
