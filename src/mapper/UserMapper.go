package mapper

import (
	"github.com/ruiborda/ecommerce-user-service/src/dto/user"
	"github.com/ruiborda/ecommerce-user-service/src/model"
	"time"
)

type UserMapper struct {}

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

func (m *UserMapper) UpdateUserRequestToUser(request *user.UpdateUserRequest, existingModel *model.User) *model.User {
	existingModel.Email = request.Email
	existingModel.FullName = request.FullName
	existingModel.RoleIds = request.RoleIds

	if request.Password != "" {
		existingModel.PasswordHash = request.Password 
	}

	existingModel.UpdatedAt = time.Now()

	return existingModel
}

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

func (m *UserMapper) UserToDeleteUserByIdResponse(userId string, success bool) *user.DeleteUserByIdResponse {
	return &user.DeleteUserByIdResponse{
		Success: success,
		Message: getDeleteMessage(userId, success),
	}
}

func getDeleteMessage(userId string, success bool) string {
	if success {
		return "User with ID " + userId + " was successfully deleted"
	}
	return "Failed to delete user with ID " + userId
}

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
