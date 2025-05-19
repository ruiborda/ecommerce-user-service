package role

type GetRolesByIdsResponse struct {
	Roles []*GetRoleByIdResponse `json:"roles"`
}
