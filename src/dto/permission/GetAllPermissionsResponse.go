package permission

type GetAllPermissionsResponse struct {
	Permissions []GetPermissionByIdResponse `json:"permissions"`
}
