package impl

import (
	"github.com/ruiborda/ecommerce-user-service/src/dto/user"
	"github.com/ruiborda/ecommerce-user-service/src/mapper"
	"github.com/ruiborda/ecommerce-user-service/src/model"
	"github.com/ruiborda/ecommerce-user-service/src/repository"
	"github.com/ruiborda/ecommerce-user-service/src/repository/impl"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	userRepository repository.UserRepository
	roleRepository repository.RoleRepository
	userMapper     *mapper.UserMapper
	roleMapper     *mapper.RoleMapper
}

func NewUserServiceImpl() *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: impl.NewUserRepositoryImpl(),
		roleRepository: impl.NewRoleRepositoryImpl(),
		userMapper:     &mapper.UserMapper{},
		roleMapper:     &mapper.RoleMapper{},
	}
}

// CreateUser crea un nuevo usuario
func (s *UserServiceImpl) CreateUser(request *user.CreateUserRequest) *user.CreateUserResponse {
	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil
	}

	// Map request to model
	userModel := s.userMapper.CreateUserRequestToUser(request)
	userModel.PasswordHash = string(passwordHash)

	// Save to database
	createdUser, err := s.userRepository.Create(userModel)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil
	}

	// Get roles for response
	var roleSlice []*model.Role
	if len(createdUser.RoleIds) > 0 {
		roleModels, err := s.roleRepository.FindByIds(createdUser.RoleIds)
		if err != nil {
			log.Printf("Error fetching roles: %v", err)
		} else {
			roleSlice = roleModels
		}
	}

	// Convert []*model.Role to []model.Role for backward compatibility
	var roles []model.Role
	for _, rolePtr := range roleSlice {
		if rolePtr != nil {
			roles = append(roles, *rolePtr)
		}
	}

	// Map model to response
	return s.userMapper.UserToCreateUserResponse(createdUser, &roles)
}

// GetUserById obtiene un usuario por su ID
func (s *UserServiceImpl) GetUserById(id string) *user.GetUserByIdResponse {
	userModel, err := s.userRepository.FindById(id)
	if err != nil {
		log.Printf("Error fetching user by ID: %v", err)
		return nil
	}

	if userModel == nil {
		return nil
	}

	// Get roles for response
	var roleSlice []*model.Role
	if len(userModel.RoleIds) > 0 {
		roleModels, err := s.roleRepository.FindByIds(userModel.RoleIds)
		if err != nil {
			log.Printf("Error fetching roles: %v", err)
		} else {
			roleSlice = roleModels
		}
	}

	// Convert []*model.Role to []model.Role for backward compatibility
	var roles []model.Role
	for _, rolePtr := range roleSlice {
		if rolePtr != nil {
			roles = append(roles, *rolePtr)
		}
	}

	return s.userMapper.UserToGetUserByIdResponse(userModel, &roles)
}

// GetUserByEmail obtiene un usuario por su email
func (s *UserServiceImpl) GetUserByEmail(email string) *user.GetUserByIdResponse {
	userModel, err := s.userRepository.FindByEmail(email)
	if err != nil {
		log.Printf("Error fetching user by email: %v", err)
		return nil
	}

	if userModel == nil {
		return nil
	}

	// Get roles for response
	var roleSlice []*model.Role
	if len(userModel.RoleIds) > 0 {
		roleModels, err := s.roleRepository.FindByIds(userModel.RoleIds)
		if err != nil {
			log.Printf("Error fetching roles: %v", err)
		} else {
			roleSlice = roleModels
		}
	}

	// Convert []*model.Role to []model.Role for backward compatibility
	var roles []model.Role
	for _, rolePtr := range roleSlice {
		if rolePtr != nil {
			roles = append(roles, *rolePtr)
		}
	}

	return s.userMapper.UserToGetUserByIdResponse(userModel, &roles)
}

// GetAllUsers obtiene todos los usuarios
func (s *UserServiceImpl) GetAllUsers() []*user.GetUserByIdResponse {
	users, err := s.userRepository.FindAll()
	if err != nil {
		log.Printf("Error fetching all users: %v", err)
		return nil
	}

	// Create map to store roles for each user
	rolesMap := make(map[string]*[]model.Role)

	// Pre-fetch roles for all users
	var allRoleIds []string
	for _, user := range users {
		allRoleIds = append(allRoleIds, user.RoleIds...)
	}

	// Deduplicate role IDs
	uniqueRoleIds := make(map[string]bool)
	var uniqueRoleIdsSlice []string
	for _, roleId := range allRoleIds {
		if !uniqueRoleIds[roleId] {
			uniqueRoleIds[roleId] = true
			uniqueRoleIdsSlice = append(uniqueRoleIdsSlice, roleId)
		}
	}

	// Fetch all needed roles in one go
	var allRolesPtr []*model.Role
	if len(uniqueRoleIdsSlice) > 0 {
		allRolesPtr, err = s.roleRepository.FindByIds(uniqueRoleIdsSlice)
		if err != nil {
			log.Printf("Error fetching roles: %v", err)
		}
	}

	// Convert []*model.Role to []model.Role for backward compatibility
	var allRoles []model.Role
	for _, rolePtr := range allRolesPtr {
		if rolePtr != nil {
			allRoles = append(allRoles, *rolePtr)
		}
	}

	// Create mapping of role ID to role
	roleById := make(map[string]model.Role)
	for _, role := range allRoles {
		roleById[role.Id] = role
	}

	// Map each user's roles
	for _, userModel := range users {
		var userRoles []model.Role
		for _, roleId := range userModel.RoleIds {
			if role, ok := roleById[roleId]; ok {
				userRoles = append(userRoles, role)
			}
		}
		rolesMap[userModel.Id] = &userRoles
	}

	// Convert users to DTOs
	var userResponses []*user.GetUserByIdResponse
	for _, userModel := range users {
		roles := rolesMap[userModel.Id]
		if roles == nil {
			emptyRoles := []model.Role{}
			roles = &emptyRoles
		}
		userDto := s.userMapper.UserToGetUserByIdResponse(userModel, roles)
		userResponses = append(userResponses, userDto)
	}

	return userResponses
}

// UpdateUserById actualiza un usuario existente
func (s *UserServiceImpl) UpdateUserById(request *user.UpdateUserRequest) *user.UpdateUserResponse {
	// First get existing user
	existingUser, err := s.userRepository.FindById(request.Id)
	if err != nil {
		log.Printf("Error fetching user to update: %v", err)
		return nil
	}

	if existingUser == nil {
		return nil
	}

	// Check if password needs to be updated
	if request.Password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return nil
		}
		request.Password = string(passwordHash)
	}

	// Map request to model
	updatedUserModel := s.userMapper.UpdateUserRequestToUser(request, existingUser)

	// Save to database
	updatedUser, err := s.userRepository.Update(updatedUserModel)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return nil
	}

	// Get roles for response
	var roleSlice []*model.Role
	if len(updatedUser.RoleIds) > 0 {
		roleModels, err := s.roleRepository.FindByIds(updatedUser.RoleIds)
		if err != nil {
			log.Printf("Error fetching roles: %v", err)
		} else {
			roleSlice = roleModels
		}
	}

	// Convert []*model.Role to []model.Role for backward compatibility
	var roles []model.Role
	for _, rolePtr := range roleSlice {
		if rolePtr != nil {
			roles = append(roles, *rolePtr)
		}
	}

	// Map model to response
	return s.userMapper.UserToUpdateUserResponse(updatedUser, &roles)
}

// DeleteUserById elimina un usuario por su ID
func (s *UserServiceImpl) DeleteUserById(id string) *user.DeleteUserByIdResponse {
	err := s.userRepository.Delete(id)
	success := err == nil

	return s.userMapper.UserToDeleteUserByIdResponse(id, success)
}

// FindAllUsersByPageAndSize obtiene usuarios paginados
func (s *UserServiceImpl) FindAllUsersByPageAndSize(page, size int) []*user.GetUserByIdResponse {
	users, err := s.userRepository.FindAllByPageAndSize(page, size)
	if err != nil {
		log.Printf("Error fetching paginated users: %v", err)
		return nil
	}

	// Create map to store roles for each user
	rolesMap := make(map[string]*[]model.Role)

	// Collect all role IDs
	var allRoleIds []string
	for _, userModel := range users {
		allRoleIds = append(allRoleIds, userModel.RoleIds...)
	}

	// Deduplicate role IDs
	uniqueRoleIds := make(map[string]bool)
	var uniqueRoleIdsSlice []string
	for _, roleId := range allRoleIds {
		if !uniqueRoleIds[roleId] {
			uniqueRoleIds[roleId] = true
			uniqueRoleIdsSlice = append(uniqueRoleIdsSlice, roleId)
		}
	}

	// Fetch all needed roles in one go
	var allRolesPtr []*model.Role
	if len(uniqueRoleIdsSlice) > 0 {
		allRolesPtr, err = s.roleRepository.FindByIds(uniqueRoleIdsSlice)
		if err != nil {
			log.Printf("Error fetching roles: %v", err)
		}
	}

	// Convert []*model.Role to []model.Role for backward compatibility
	var allRoles []model.Role
	for _, rolePtr := range allRolesPtr {
		if rolePtr != nil {
			allRoles = append(allRoles, *rolePtr)
		}
	}

	// Create mapping of role ID to role
	roleById := make(map[string]model.Role)
	for _, role := range allRoles {
		roleById[role.Id] = role
	}

	// Map each user's roles
	for _, userModel := range users {
		var userRoles []model.Role
		for _, roleId := range userModel.RoleIds {
			if role, ok := roleById[roleId]; ok {
				userRoles = append(userRoles, role)
			}
		}
		rolesMap[userModel.Id] = &userRoles
	}

	// Convert users to DTOs
	var userResponses []*user.GetUserByIdResponse
	for _, userModel := range users {
		roles := rolesMap[userModel.Id]
		if roles == nil {
			emptyRoles := []model.Role{}
			roles = &emptyRoles
		}
		userDto := s.userMapper.UserToGetUserByIdResponse(userModel, roles)
		userResponses = append(userResponses, userDto)
	}

	return userResponses
}

// CountAllUsers cuenta el número total de usuarios
func (s *UserServiceImpl) CountAllUsers() int64 {
	count, err := s.userRepository.Count()
	if err != nil {
		log.Printf("Error counting users: %v", err)
		return 0
	}

	return count
}

// GetUsersByIds obtiene múltiples usuarios por sus IDs
func (s *UserServiceImpl) GetUsersByIds(ids []string) []*user.GetUserByIdResponse {
	users, err := s.userRepository.FindByIds(ids)
	if err != nil {
		log.Printf("Error fetching users by IDs: %v", err)
		return nil
	}

	// Create map to store roles for each user
	rolesMap := make(map[string]*[]model.Role)

	// Collect all role IDs
	var allRoleIds []string
	for _, userModel := range users {
		allRoleIds = append(allRoleIds, userModel.RoleIds...)
	}

	// Deduplicate role IDs
	uniqueRoleIds := make(map[string]bool)
	var uniqueRoleIdsSlice []string
	for _, roleId := range allRoleIds {
		if !uniqueRoleIds[roleId] {
			uniqueRoleIds[roleId] = true
			uniqueRoleIdsSlice = append(uniqueRoleIdsSlice, roleId)
		}
	}

	// Fetch all needed roles in one go
	var allRolesPtr []*model.Role
	if len(uniqueRoleIdsSlice) > 0 {
		allRolesPtr, err = s.roleRepository.FindByIds(uniqueRoleIdsSlice)
		if err != nil {
			log.Printf("Error fetching roles: %v", err)
		}
	}

	// Convert []*model.Role to []model.Role for backward compatibility
	var allRoles []model.Role
	for _, rolePtr := range allRolesPtr {
		if rolePtr != nil {
			allRoles = append(allRoles, *rolePtr)
		}
	}

	// Create mapping of role ID to role
	roleById := make(map[string]model.Role)
	for _, role := range allRoles {
		roleById[role.Id] = role
	}

	// Map each user's roles
	for _, userModel := range users {
		var userRoles []model.Role
		for _, roleId := range userModel.RoleIds {
			if role, ok := roleById[roleId]; ok {
				userRoles = append(userRoles, role)
			}
		}
		rolesMap[userModel.Id] = &userRoles
	}

	// Convert the users to individual responses
	var userResponses []*user.GetUserByIdResponse
	for _, userModel := range users {
		roles := rolesMap[userModel.Id]
		if roles == nil {
			emptyRoles := []model.Role{}
			roles = &emptyRoles
		}
		userDto := s.userMapper.UserToGetUserByIdResponse(userModel, roles)
		userResponses = append(userResponses, userDto)
	}

	return userResponses
}

// Helper methods for authentication
// These methods are used by AuthServiceImpl

// GetRolesForUser returns the roles for a given user
func (s *UserServiceImpl) GetRolesForUser(user *model.User) *[]model.Role {
	if len(user.RoleIds) == 0 {
		emptyRoles := []model.Role{}
		return &emptyRoles
	}

	roleModelsPtr, err := s.roleRepository.FindByIds(user.RoleIds)
	if err != nil {
		log.Printf("Error fetching roles for user: %v", err)
		emptyRoles := []model.Role{}
		return &emptyRoles
	}

	// Convert []*model.Role to []model.Role
	var roleModels []model.Role
	for _, rolePtr := range roleModelsPtr {
		if rolePtr != nil {
			roleModels = append(roleModels, *rolePtr)
		}
	}

	return &roleModels
}

// VerifyPassword checks if a password matches the hash
func (s *UserServiceImpl) VerifyPassword(user *model.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}

// SetPasswordHash hashes and sets a password for a user
func (s *UserServiceImpl) SetPasswordHash(user *model.User, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(passwordHash)
	return nil
}

// CreateUserWithRoleAndAuth creates a new user with specified roles
func (s *UserServiceImpl) CreateUserWithRoleAndAuth(email, fullName, password string, roleIds []string) (*model.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepository.FindByEmail(email)
	if err == nil && existingUser != nil {
		return nil, nil // User already exists
	}

	// Create new user
	now := time.Now()
	userModel := &model.User{
		Email:                  email,
		FullName:               fullName,
		CreatedAt:              now,
		UpdatedAt:              now,
		RoleIds:                roleIds,
		FavoriteNewsArticleIds: []string{},
	}

	// Set password if provided
	if password != "" {
		err = s.SetPasswordHash(userModel, password)
		if err != nil {
			return nil, err
		}
	}

	// Save to database
	createdUser, err := s.userRepository.Create(userModel)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
