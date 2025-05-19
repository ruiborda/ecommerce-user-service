package model

import (
	"time"
)

type User struct {
	Id    string `json:"id" firestore:"id,omitempty"`
	Email string `json:"email" firestore:"email,omitempty"`
	// usar bycrypt to hash password
	PasswordHash string `json:"passwordHash" firestore:"passwordHash,omitempty"`
	FullName     string `json:"fullName" firestore:"fullName,omitempty"`
	ImageFileKey string `json:"imageFileKey" firestore:"imageFileKey,omitempty"`
	// google picture url of login
	PictureUrl string    `json:"pictureUrl" firestore:"pictureUrl,omitempty"`
	CreatedAt  time.Time `json:"createdAt" firestore:"createdAt,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt" firestore:"updatedAt,omitempty"`
	RoleIds    []string  `json:"roleIds" firestore:"roleIds,omitempty"`
	// IDs of favorite news articles
	FavoriteNewsArticleIds []string `json:"favoriteNewsArticleIds" firestore:"favoriteNewsArticleIds,omitempty"`
}
