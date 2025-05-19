package impl

import (
	"UserService/src/dto/user"
	"UserService/src/mapper"
	"UserService/src/model"
	"UserService/src/repository"
	"UserService/src/repository/impl"
	"log"

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

// CreateUser creates a new user
func (s *UserServiceImpl) CreateUser(request *user.CreateUserRequest) *user.CreateUserResponse {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil
	}

	// Map request to model
	userModel := s.userMapper.CreateUserRequestToUser(request)
	userModel.PasswordHash = string(hashedPassword)

	// Save to repository
	createdUser, err := s.userRepository.Create(userModel)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil
	}

	// Get roles if roleIds are provided
	var roles *[]model.Role
	if len(createdUser.RoleIds) > 0 {
		rolesData, err := s.roleRepository.FindByIds(createdUser.RoleIds)
		if err != nil {
			log.Printf("Error fetching roles: %v", err)
		} else {
			roles = &rolesData
		}
	}

	// Map created user to response
	return s.userMapper.UserToCreateUserResponse(createdUser, roles)
}

// GetUserById retrieves a user by their ID
func (s *UserServiceImpl) GetUserById(id string) *user.GetUserByIdResponse {
	userModel, err := s.userRepository.FindById(id)
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		return nil
	}

	if userModel == nil {
		return nil
	}

	// Get roles if roleIds are provided
	var roles *[]model.Role
	if len(userModel.RoleIds) > 0 {
		rolesData, err := s.roleRepository.FindByIds(userModel.RoleIds)
		if err != nil {
			log.Printf("Error fetching roles: %v", err)
		} else {
			roles = &rolesData
		}
	}

	return s.userMapper.UserToGetUserByIdResponse(userModel, roles)
}

// GetUserByEmail retrieves a user by their email
func (s *UserServiceImpl) GetUserByEmail(email string) *user.GetUserByIdResponse {
	userModel, err := s.userRepository.FindByEmail(email)
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		return nil
	}

	if userModel == nil {
		return nil
	}

	// Get roles if roleIds are provided
	var roles *[]model.Role
	if len(userModel.RoleIds) > 0 {
		rolesData, err := s.roleRepository.FindByIds(userModel.RoleIds)
		if err != nil {
			log.Printf("Error fetching roles: %v", err)
		} else {
			roles = &rolesData
		}
	}

	return s.userMapper.UserToGetUserByIdResponse(userModel, roles)
}

// GetAllUsers retrieves all users
func (s *UserServiceImpl) GetAllUsers() []*user.GetUserByIdResponse {
	users, err := s.userRepository.FindAll()
	if err != nil {
		log.Printf("Error getting all users: %v", err)
		return []*user.GetUserByIdResponse{} // Devolver un array vacío en lugar de nil
	}

	// Si users es nil, devolver un array vacío
	if users == nil {
		return []*user.GetUserByIdResponse{}
	}

	// Create a map to store roles for each user
	rolesMap := make(map[string]*[]model.Role)

	// Collect all role IDs from all users
	var allRoleIds []string
	for _, user := range users {
		for _, roleId := range user.RoleIds {
			allRoleIds = append(allRoleIds, roleId)
		}
	}

	// Fetch all roles at once (to minimize database calls)
	if len(allRoleIds) > 0 {
		allRoles, err := s.roleRepository.FindByIds(allRoleIds)
		if err == nil {
			// Create a map of role ID to role for easy lookup
			roleById := make(map[string]model.Role)
			for _, role := range allRoles {
				roleById[role.Id] = role
			}

			// Associate roles with each user
			for _, user := range users {
				var userRoles []model.Role
				for _, roleId := range user.RoleIds {
					if role, exists := roleById[roleId]; exists {
						userRoles = append(userRoles, role)
					}
				}
				if len(userRoles) > 0 {
					rolesMap[user.Id] = &userRoles
				}
			}
		} else {
			log.Printf("Error fetching roles: %v", err)
		}
	}

	// Usar el nuevo método de mapper
	response := s.userMapper.UsersToGetUsersByIdsResponse(users, rolesMap)

	// Si response.Users es nil, devolver un array vacío
	if response.Users == nil {
		return []*user.GetUserByIdResponse{}
	}

	return response.Users
}

// UpdateUserById updates an existing user
func (s *UserServiceImpl) UpdateUserById(request *user.UpdateUserRequest) *user.UpdateUserResponse {
	// First get the existing user
	existingUser, err := s.userRepository.FindById(request.Id)
	if err != nil {
		log.Printf("Error finding user to update: %v", err)
		return nil
	}

	if existingUser == nil {
		log.Printf("User not found with ID: %s", request.Id)
		return nil
	}

	// Hash the password if provided
	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return nil
		}
		request.Password = string(hashedPassword)
	}

	// Update the user model with request data
	updatedUserModel := s.userMapper.UpdateUserRequestToUser(request, existingUser)

	// Save the updated user
	savedUser, err := s.userRepository.Update(updatedUserModel)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return nil
	}

	// Get roles if roleIds are provided
	var roles *[]model.Role
	if len(savedUser.RoleIds) > 0 {
		rolesData, err := s.roleRepository.FindByIds(savedUser.RoleIds)
		if err != nil {
			log.Printf("Error fetching roles: %v", err)
		} else {
			roles = &rolesData
		}
	}

	return s.userMapper.UserToUpdateUserResponse(savedUser, roles)
}

// DeleteUserById deletes a user by their ID
func (s *UserServiceImpl) DeleteUserById(id string) *user.DeleteUserByIdResponse {
	err := s.userRepository.Delete(id)
	success := err == nil

	// Create the response using the mapper regardless of success or failure
	response := s.userMapper.UserToDeleteUserByIdResponse(id, success)

	if err != nil {
		log.Printf("Error deleting user: %v", err)
		// Return response with success=false
	}

	return response
}

// FindAllUsersByPageAndSize retrieves users with pagination
func (s *UserServiceImpl) FindAllUsersByPageAndSize(page, size int) []*user.GetUserByIdResponse {
	users, err := s.userRepository.FindAllByPageAndSize(page, size)
	if err != nil {
		log.Printf("Error getting paginated users: %v", err)
		return []*user.GetUserByIdResponse{} // Devolver un array vacío en lugar de nil
	}

	// Si users es nil, devolver un array vacío
	if users == nil {
		return []*user.GetUserByIdResponse{}
	}

	// Create a map to store roles for each user
	rolesMap := make(map[string]*[]model.Role)

	// Collect all role IDs from all users
	var allRoleIds []string
	for _, user := range users {
		for _, roleId := range user.RoleIds {
			allRoleIds = append(allRoleIds, roleId)
		}
	}

	// Fetch all roles at once (to minimize database calls)
	if len(allRoleIds) > 0 {
		allRoles, err := s.roleRepository.FindByIds(allRoleIds)
		if err == nil {
			// Create a map of role ID to role for easy lookup
			roleById := make(map[string]model.Role)
			for _, role := range allRoles {
				roleById[role.Id] = role
			}

			// Associate roles with each user
			for _, user := range users {
				var userRoles []model.Role
				for _, roleId := range user.RoleIds {
					if role, exists := roleById[roleId]; exists {
						userRoles = append(userRoles, role)
					}
				}
				if len(userRoles) > 0 {
					rolesMap[user.Id] = &userRoles
				}
			}
		} else {
			log.Printf("Error fetching roles: %v", err)
		}
	}

	// Usar el nuevo método de mapper
	response := s.userMapper.UsersToGetUsersByIdsResponse(users, rolesMap)

	// Si response.Users es nil, devolver un array vacío
	if response.Users == nil {
		return []*user.GetUserByIdResponse{}
	}

	return response.Users
}

// CountAllUsers gets the total count of users
func (s *UserServiceImpl) CountAllUsers() int64 {
	count, err := s.userRepository.Count()
	if err != nil {
		log.Printf("Error counting users: %v", err)
		return 0
	}
	return count
}

// GetUsersByIds retrieves users by a slice of IDs
func (s *UserServiceImpl) GetUsersByIds(ids []string) []*user.GetUserByIdResponse {
	users, err := s.userRepository.FindByIds(ids)
	if err != nil {
		log.Printf("Error getting users by IDs: %v", err)
		return []*user.GetUserByIdResponse{} // Devolver un array vacío en lugar de nil
	}

	// Si users es nil, devolver un array vacío
	if users == nil {
		return []*user.GetUserByIdResponse{}
	}

	// Create a map to store roles for each user
	rolesMap := make(map[string]*[]model.Role)

	// Collect all role IDs from all users
	var allRoleIds []string
	for _, user := range users {
		for _, roleId := range user.RoleIds {
			allRoleIds = append(allRoleIds, roleId)
		}
	}

	// Fetch all roles at once (to minimize database calls)
	if len(allRoleIds) > 0 {
		allRoles, err := s.roleRepository.FindByIds(allRoleIds)
		if err == nil {
			// Create a map of role ID to role for easy lookup
			roleById := make(map[string]model.Role)
			for _, role := range allRoles {
				roleById[role.Id] = role
			}

			// Associate roles with each user
			for _, user := range users {
				var userRoles []model.Role
				for _, roleId := range user.RoleIds {
					if role, exists := roleById[roleId]; exists {
						userRoles = append(userRoles, role)
					}
				}
				if len(userRoles) > 0 {
					rolesMap[user.Id] = &userRoles
				}
			}
		} else {
			log.Printf("Error fetching roles: %v", err)
		}
	}

	// Usar el nuevo método de mapper
	response := s.userMapper.UsersToGetUsersByIdsResponse(users, rolesMap)

	// Si response.Users es nil, devolver un array vacío
	if response.Users == nil {
		return []*user.GetUserByIdResponse{}
	}

	return response.Users
}
