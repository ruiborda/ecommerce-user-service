package auth

type LoginWithAnyResponse struct {
	Id           string `json:"id"`
	FullName     string `json:"fullName"`
	GivenName    string `json:"givenName"`
	FamilyName   string `json:"familyName"`
	ProfileImage string `json:"profileImage"`
	Email        string `json:"email"`
	Jwt          string `json:"jwt"`
}
