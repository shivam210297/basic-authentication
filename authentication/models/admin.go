package models

type GetToken struct {
	Token string `json:"token"`
}

type UserContext struct {
	ID    string `db:"id"`
	Email string `db:"email"`
	Name  string `db:"name"`
}

type TokenDetails struct {
	Token      string `json:"token"db:"invite_token"`
	IsArchived bool   `json:"isArchived"db:"is_archived"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
