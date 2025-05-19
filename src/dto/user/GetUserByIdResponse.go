package user

import (
	"UserService/src/model"
	"time"
)

type GetUserByIdResponse struct {
	Id                    string        `json:"id"`
	Email                 string        `json:"email"`
	FullName              string        `json:"fullName"`
	ImageFileKey          string        `json:"imageFileKey,omitempty"`
	PictureUrl            string        `json:"pictureUrl,omitempty"`
	CreatedAt             time.Time     `json:"createdAt"`
	UpdatedAt             time.Time     `json:"updatedAt"`
	Roles                 *[]model.Role `json:"roles,omitempty"`
	FavoriteNewsArticleIds []string      `json:"favoriteNewsArticleIds,omitempty"`
}
