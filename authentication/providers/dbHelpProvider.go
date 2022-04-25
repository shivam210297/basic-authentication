package providers

import "Assignment/models"

type DBHelpProvider interface {
	LogInUserUsingEmail(loginReq models.LoginRequest) (userID string, message string, err error)
	AddToken(inviteTokenDetail models.TokenDetail) error
	DeleteToken(tokenID string) error
	VerifyToken(token string) (bool, error)
	PopulateCache() ([]models.TokenDetail, error)
	GetTokens() ([]models.TokenDetails, error)
}
