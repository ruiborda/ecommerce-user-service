package impl

import (
	"UserService/src/database"
	"UserService/src/model"
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type UserRepositoryImpl struct {
	collectionName string
}

func NewUserRepositoryImpl() *UserRepositoryImpl {
	return &UserRepositoryImpl{
		collectionName: "users",
	}
}

func (r *UserRepositoryImpl) Create(user *model.User) (*model.User, error) {
	ctx := context.Background()
	client := database.GetFirestoreClient()

	if user.Id == "" {
		user.Id = uuid.New().String()
		// Guardar el documento con el ID generado
		_, err := client.Collection(r.collectionName).Doc(user.Id).Set(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to create user with generated UUID: %v", err)
		}
	} else {
		// Verificar que el ID sea un UUID válido
		if _, err := uuid.Parse(user.Id); err != nil {
			return nil, fmt.Errorf("invalid UUID format for user ID: %v", err)
		}

		_, err := client.Collection(r.collectionName).Doc(user.Id).Set(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to create user with provided ID: %v", err)
		}
	}

	return user, nil
}

func (r *UserRepositoryImpl) FindById(id string) (*model.User, error) {
	ctx := context.Background()
	client := database.GetFirestoreClient()

	docSnap, err := client.Collection(r.collectionName).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	var user model.User
	if err := docSnap.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to convert document to user: %v", err)
	}

	// Ensure the ID is set
	user.Id = docSnap.Ref.ID

	return &user, nil
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*model.User, error) {
	ctx := context.Background()
	client := database.GetFirestoreClient()

	query := client.Collection(r.collectionName).Where("email", "==", email).Limit(1)
	iter := query.Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user by email: %v", err)
	}

	var user model.User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to convert document to user: %v", err)
	}

	// Ensure the ID is set
	user.Id = doc.Ref.ID

	return &user, nil
}

func (r *UserRepositoryImpl) FindAll() ([]*model.User, error) {
	ctx := context.Background()
	client := database.GetFirestoreClient()

	iter := client.Collection(r.collectionName).Documents(ctx)
	defer iter.Stop()

	var users []*model.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate users: %v", err)
		}

		var user model.User
		if err := doc.DataTo(&user); err != nil {
			return nil, fmt.Errorf("failed to convert document to user: %v", err)
		}

		// Ensure the ID is set
		user.Id = doc.Ref.ID
		users = append(users, &user)
	}

	return users, nil
}

func (r *UserRepositoryImpl) Update(user *model.User) (*model.User, error) {
	ctx := context.Background()
	client := database.GetFirestoreClient()

	// Update the timestamp
	user.UpdatedAt = time.Now()

	_, err := client.Collection(r.collectionName).Doc(user.Id).Set(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return user, nil
}

func (r *UserRepositoryImpl) Delete(id string) error {
	ctx := context.Background()
	client := database.GetFirestoreClient()

	_, err := client.Collection(r.collectionName).Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}

func (r *UserRepositoryImpl) FindAllByPageAndSize(page, size int) ([]*model.User, error) {
	ctx := context.Background()
	client := database.GetFirestoreClient()

	// Calculate offset
	offset := page * size

	// Get all documents and apply pagination in memory
	// Note: Firestore doesn't natively support offset-based pagination
	iter := client.Collection(r.collectionName).OrderBy("createdAt", firestore.Desc).Documents(ctx)
	defer iter.Stop()

	var users []*model.User
	index := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate users: %v", err)
		}

		// Skip documents before offset
		if index < offset {
			index++
			continue
		}

		// Break if we've collected enough documents
		if len(users) >= size {
			break
		}

		var user model.User
		if err := doc.DataTo(&user); err != nil {
			return nil, fmt.Errorf("failed to convert document to user: %v", err)
		}

		// Ensure the ID is set
		user.Id = doc.Ref.ID
		users = append(users, &user)
		index++
	}

	return users, nil
}

func (r *UserRepositoryImpl) Count() (int64, error) {
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
			return 0, fmt.Errorf("failed to count users: %v", err)
		}
		count++
	}

	return count, nil
}

func (r *UserRepositoryImpl) FindByIds(ids []string) ([]*model.User, error) {
	ctx := context.Background()
	client := database.GetFirestoreClient()

	var users []*model.User

	// Firestore doesn't support "IN" queries for document IDs
	// Need to fetch documents one by one
	for _, id := range ids {
		doc, err := client.Collection(r.collectionName).Doc(id).Get(ctx)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				// Skip non-existent documents
				continue
			}
			return nil, fmt.Errorf("failed to get user with ID %s: %v", id, err)
		}

		var user model.User
		if err := doc.DataTo(&user); err != nil {
			return nil, fmt.Errorf("failed to convert document to user: %v", err)
		}

		// Ensure the ID is set
		user.Id = doc.Ref.ID
		users = append(users, &user)
	}

	return users, nil
}
