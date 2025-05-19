package role

type CreateRoleRequest struct {
	Code        string `json:"code"`
	Permissions []int  `json:"permissions"`
}
