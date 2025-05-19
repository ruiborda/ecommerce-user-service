package user

type GetUsersByIdsResponse struct {
	Users []*GetUserByIdResponse `json:"users"`
}