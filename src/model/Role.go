package model

type Role struct {
	Id          string        `json:"id" firestore:"id,omitempty"`
	Code        string        `json:"code" firestore:"code,omitempty"`
	Permissions *[]Permission `json:"permissions" firestore:"permissions,omitempty"`
}
