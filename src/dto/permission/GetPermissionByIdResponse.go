package permission

type GetPermissionByIdResponse struct {
	Id          int    `json:"id"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Name        string `json:"name"`
	Description string `json:"description"`
}