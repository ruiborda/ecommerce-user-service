package role

import "UserService/src/model"

type UpdateRoleResponse struct {
	Id          string              `json:"id"`
	Code        string              `json:"code"`
	Permissions *[]model.Permission `json:"permissions"`
}
