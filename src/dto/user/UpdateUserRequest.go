package user

type UpdateUserRequest struct {
	Id       string   `json:"id"`
	Email    string   `json:"email"`
	Password string   `json:"password,omitempty"`
	FullName string   `json:"fullName"`
	RoleIds  []string `json:"roleIds"`
}
