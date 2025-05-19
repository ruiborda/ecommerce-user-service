// filepath: /home/rui/ecommerce/UserService/src/repository/impl/RoleRepositoryImpl.go
package impl

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ruiborda/ecommerce-user-service/src/database"
	"github.com/ruiborda/ecommerce-user-service/src/model"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RoleRepositoryImpl struct {
	collectionName string
}

func NewRoleRepositoryImpl() *RoleRepositoryImpl {
	return &RoleRepositoryImpl{
		collectionName: "roles",
	}
}

func (r *RoleRepositoryImpl) Create(role *model.Role) (*model.Role, error) {
	ctx := context.Background()
	client := database.GetFirestoreClient()

	if role.Id == "" {
		// Generar UUID para nuevos roles
		role.Id = uuid.New().String()

		// Guardar el documento con el ID generado
		_, err := client.Collection(r.collectionName).Doc(role.Id).Set(ctx, role)
		if err != nil {
			return nil, fmt.Errorf("failed to create role with generated UUID: %v", err)
		}
	} else {
		// Verificar que el ID sea un UUID válido
		if _, err := uuid.Parse(role.Id); err != nil {
			return nil, fmt.Errorf("invalid UUID format for role ID: %v", err)
		}

		_, err := client.Collection(r.collectionName).Doc(role.Id).Set(ctx, role)
		if err != nil {
			return nil, fmt.Errorf("failed to create role with provided ID: %v", err)
		}
	}

	return role, nil
}

func (r *RoleRepositoryImpl) FindById(id string) (*model.Role, error) {
	ctx := context.Background()
	client := database.GetFirestoreClient()

	docSnap, err := client.Collection(r.collectionName).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get role: %v", err)
	}

	var role model.Role
	if err := docSnap.DataTo(&role); err != nil {
		return nil, fmt.Errorf("failed to convert document to role: %v", err)
	}

	// Ensure the ID is set
	role.Id = docSnap.Ref.ID

	return &role, nil
}

func (r *RoleRepositoryImpl) FindByCode(code string) (*model.Role, error) {
	ctx := context.Background()
	client := database.GetFirestoreClient()

	query := client.Collection(r.collectionName).Where("code", "==", code).Limit(1)
	iter := query.Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query role by code: %v", err)
	}

	var role model.Role
	if err := doc.DataTo(&role); err != nil {
		return nil, fmt.Errorf("failed to convert document to role: %v", err)
	}

	// Ensure the ID is set
	role.Id = doc.Ref.ID

	return &role, nil
}

func (r *RoleRepositoryImpl) FindAll() ([]*model.Role, error) {
	ctx := context.Background()
	client := database.GetFirestoreClient()

	iter := client.Collection(r.collectionName).Documents(ctx)
	defer iter.Stop()

	var roles []*model.Role
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate roles: %v", err)
		}

		var role model.Role
		if err := doc.DataTo(&role); err != nil {
			return nil, fmt.Errorf("failed to convert document to role: %v", err)
		}

		// Ensure the ID is set
		role.Id = doc.Ref.ID
		roles = append(roles, &role)
	}

	return roles, nil
}

func (r *RoleRepositoryImpl) Update(role *model.Role) (*model.Role, error) {
	ctx := context.Background()
	client := database.GetFirestoreClient()

	_, err := client.Collection(r.collectionName).Doc(role.Id).Set(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("failed to update role: %v", err)
	}

	return role, nil
}

func (r *RoleRepositoryImpl) Delete(id string) error {
	ctx := context.Background()
	client := database.GetFirestoreClient()

	_, err := client.Collection(r.collectionName).Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete role: %v", err)
	}

	return nil
}

func (r *RoleRepositoryImpl) FindAllByPageAndSize(page, size int) ([]*model.Role, error) {
	ctx := context.Background()
	client := database.GetFirestoreClient()

	// Calculate offset
	offset := page * size

	// Get all documents and apply pagination in memory
	// Note: Firestore doesn't natively support offset-based pagination
	iter := client.Collection(r.collectionName).OrderBy("code", firestore.Asc).Documents(ctx)
	defer iter.Stop()

	var roles []*model.Role
	index := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate roles: %v", err)
		}

		// Skip documents before offset
		if index < offset {
			index++
			continue
		}

		// Break if we've collected enough documents
		if len(roles) >= size {
			break
		}

		var role model.Role
		if err := doc.DataTo(&role); err != nil {
			return nil, fmt.Errorf("failed to convert document to role: %v", err)
		}

		// Ensure the ID is set
		role.Id = doc.Ref.ID
		roles = append(roles, &role)
		index++
	}

	return roles, nil
}

func (r *RoleRepositoryImpl) Count() (int64, error) {
	ctx := context.Background()
	client := database.GetFirestoreClient()

	// Firestore doesn't provide a direct count operation
	// We need to iterate through all documents
	iter := client.Collection(r.collectionName).Documents(ctx)
	defer iter.Stop()

	var count int64
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("failed to count roles: %v", err)
		}
		count++
	}

	return count, nil
}

func (r *RoleRepositoryImpl) FindByIds(ids []string) ([]*model.Role, error) {
	ctx := context.Background()
	client := database.GetFirestoreClient()

	var roles []*model.Role

	// Fetch each role by ID
	for _, id := range ids {
		doc, err := client.Collection(r.collectionName).Doc(id).Get(ctx)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				// Skip non-existent documents
				continue
			}
			return nil, fmt.Errorf("failed to get role with ID %s: %v", id, err)
		}

		var role model.Role
		if err := doc.DataTo(&role); err != nil {
			return nil, fmt.Errorf("failed to convert document to role: %v", err)
		}

		// Ensure the ID is set
		role.Id = doc.Ref.ID
		roles = append(roles, &role)
	}

	return roles, nil
}
