package user

type CreateUserRequest struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	FullName string   `json:"fullName"`
	RoleIds  []string `json:"roleIds"`
}
