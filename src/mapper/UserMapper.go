package mapper

import (
	"UserService/src/dto/user"
	"UserService/src/model"
	"time"
)

type UserMapper struct {
}

// CreateUserRequestToUser convierte un CreateUserRequest a un modelo User
func (m *UserMapper) CreateUserRequestToUser(request *user.CreateUserRequest) *model.User {
	return &model.User{
		Email:                  request.Email,
		FullName:               request.FullName,
		RoleIds:                request.RoleIds,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
		FavoriteNewsArticleIds: []string{},
	}
}

// UserToCreateUserResponse convierte un modelo User a un CreateUserResponse
func (m *UserMapper) UserToCreateUserResponse(model *model.User, roles *[]model.Role) *user.CreateUserResponse {
	return &user.CreateUserResponse{
		Id:           model.Id,
		Email:        model.Email,
		FullName:     model.FullName,
		ImageFileKey: model.ImageFileKey,
		PictureUrl:   model.PictureUrl,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
		Roles:        roles,
	}
}

// GetUserByIdRequestToUser (no es necesario, ya que es solo un ID)

// UserToGetUserByIdResponse convierte un modelo User a un GetUserByIdResponse
func (m *UserMapper) UserToGetUserByIdResponse(model *model.User, roles *[]model.Role) *user.GetUserByIdResponse {
	return &user.GetUserByIdResponse{
		Id:                     model.Id,
		Email:                  model.Email,
		FullName:               model.FullName,
		ImageFileKey:           model.ImageFileKey,
		PictureUrl:             model.PictureUrl,
		CreatedAt:              model.CreatedAt,
		UpdatedAt:              model.UpdatedAt,
		Roles:                  roles,
		FavoriteNewsArticleIds: model.FavoriteNewsArticleIds,
	}
}

// UpdateUserRequestToUser convierte un UpdateUserRequest a un modelo User
func (m *UserMapper) UpdateUserRequestToUser(request *user.UpdateUserRequest, existingModel *model.User) *model.User {
	// Mantener los datos que no se deben modificar
	existingModel.Email = request.Email
	existingModel.FullName = request.FullName
	existingModel.RoleIds = request.RoleIds

	// Si hay una nueva contraseña, actualizar el hash de la contraseña
	if request.Password != "" {
		existingModel.PasswordHash = request.Password // Ya está hasheada en el servicio
	}

	// Actualizar timestamp
	existingModel.UpdatedAt = time.Now()

	return existingModel
}

// UserToUpdateUserResponse convierte un modelo User a un UpdateUserResponse
func (m *UserMapper) UserToUpdateUserResponse(model *model.User, roles *[]model.Role) *user.UpdateUserResponse {
	return &user.UpdateUserResponse{
		Id:                     model.Id,
		Email:                  model.Email,
		FullName:               model.FullName,
		ImageFileKey:           model.ImageFileKey,
		PictureUrl:             model.PictureUrl,
		CreatedAt:              model.CreatedAt,
		UpdatedAt:              model.UpdatedAt,
		Roles:                  roles,
		FavoriteNewsArticleIds: model.FavoriteNewsArticleIds,
	}
}

// DeleteUserByIdRequestToUser (no es necesario, ya que es solo un ID)

// UserToDeleteUserByIdResponse convierte un resultado de eliminación a un DeleteUserByIdResponse
func (m *UserMapper) UserToDeleteUserByIdResponse(userId string, success bool) *user.DeleteUserByIdResponse {
	return &user.DeleteUserByIdResponse{
		Success: success,
		Message: getDeleteMessage(userId, success),
	}
}

// Función auxiliar para crear un mensaje adecuado
func getDeleteMessage(userId string, success bool) string {
	if success {
		return "User with ID " + userId + " was successfully deleted"
	}
	return "Failed to delete user with ID " + userId
}

// GetUsersByIdsRequestToUsers (no es necesario, ya que es solo una lista de IDs)

// UsersToGetUsersByIdsResponse convierte una lista de modelos User a una lista de GetUserByIdResponse
func (m *UserMapper) UsersToGetUsersByIdsResponse(models []*model.User, rolesMap map[string]*[]model.Role) *user.GetUsersByIdsResponse {
	var userResponses []*user.GetUserByIdResponse

	for _, model := range models {
		roles, exists := rolesMap[model.Id]

		response := &user.GetUserByIdResponse{
			Id:                     model.Id,
			Email:                  model.Email,
			FullName:               model.FullName,
			ImageFileKey:           model.ImageFileKey,
			PictureUrl:             model.PictureUrl,
			CreatedAt:              model.CreatedAt,
			UpdatedAt:              model.UpdatedAt,
			FavoriteNewsArticleIds: model.FavoriteNewsArticleIds,
		}

		if exists {
			response.Roles = roles
		}

		userResponses = append(userResponses, response)
	}

	return &user.GetUsersByIdsResponse{
		Users: userResponses,
	}
}
