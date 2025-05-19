package role

type UpdateRoleRequest struct {
	Id          string `json:"id"`
	Code        string `json:"code"`
	Permissions []int  `json:"permissions"`
}
