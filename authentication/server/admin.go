package server

import (
	"Assignment/models"
	"Assignment/utils"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

// login - api for the admin login
// user login is based on user email and password
// checking some initial validation on payload
// Using JWT token for authentication
func (srv *Server) login(resp http.ResponseWriter, req *http.Request) {
	var loginRequest models.LoginRequest
	err := json.NewDecoder(req.Body).Decode(&loginRequest)
	if err != nil {
		logrus.Errorf("login: error in parsing json: %v", err)
		return
	}

	if loginRequest.Email == "" {
		utils.EncodeJSONBody(resp, http.StatusBadRequest, map[string]interface{}{
			"message": "email cannot be empty",
		})
		return
	}

	if loginRequest.Password == "" {
		utils.EncodeJSONBody(resp, http.StatusBadRequest, map[string]interface{}{
			"message": "password cannot be empty",
		})
		return
	}

	userID, message, err := srv.DBHelper.LogInUserUsingEmail(loginRequest)
	if err != nil {
		logrus.Errorf("login: error in validating user: message: %v: err : %v", message, err)
		utils.EncodeJSONBody(resp, http.StatusInternalServerError, map[string]interface{}{
			"message": message,
		})
		return
	}

	// Create a struct that will be encoded to a JWT.
	// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
	type Claims struct {
		UserID string `json:"userID"`
		jwt.StandardClaims
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	var signingKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		utils.EncodeJSONBody(resp, http.StatusInternalServerError, map[string]interface{}{
			"message": "error while validating user",
		})
		return
	}

	utils.EncodeJSON200Body(resp, map[string]interface{}{
		"token": tokenString,
	})
}

// GenerateToken - this code id used for generating invite code
// created a logic for generating a alphanumeric invite code
// inserting the code in database
// if insert was successful then updating the cache
func (srv Server) GenerateToken() http.HandlerFunc {
	return func(resp http.ResponseWriter, _ *http.Request) {
		inviteToken, err := utils.CreateInviteCode(srv.cache)
		if err != nil {
			logrus.Errorf("GenerateToken: error in generating new token")
			utils.EncodeJSONBody(resp, http.StatusInternalServerError, map[string]interface{}{
				"message": "maximum limit exceeded while generating token ",
			})
			return
		}

		tokenValues := models.TokenDetail{
			Token:          inviteToken,
			Count:          0,
			ExpirationTime: time.Now().AddDate(0, 0, 7),
		}

		err = srv.DBHelper.AddToken(tokenValues)
		if err != nil {
			utils.EncodeJSONBody(resp, http.StatusInternalServerError, map[string]interface{}{
				"message": "error in inserting invite code",
			})
			return
		}

		err = srv.CacheProvider.Set(inviteToken, tokenValues)
		if err != nil {
			logrus.Errorf("GenerateToken: error in setting token in cache: %v", err)
		}

		utils.EncodeJSON200Body(resp, map[string]interface{}{
			"inviteCode": inviteToken,
		})
	}
}

// disableToken - api to disable to the specified token
// taking invite token as query param
// then updating the database and cache
// returning a success message
func (srv *Server) disableToken(resp http.ResponseWriter, req *http.Request) {
	tokenID := chi.URLParam(req, "tokenID")

	err := srv.DBHelper.DeleteToken(tokenID)
	if err != nil {
		utils.EncodeJSONBody(resp, http.StatusInternalServerError, map[string]interface{}{
			"message": "error in inserting invite code",
		})
		return
	}

	srv.CacheProvider.Delete(tokenID)

	utils.EncodeJSON200Body(resp, map[string]interface{}{
		"message": "disabled successfully",
	})
}

// GetTokens - api to get the list of tokens generated
// the response contains the token and a boolean which represent whether the token is active or not
func (srv *Server) GetTokens(resp http.ResponseWriter, req *http.Request) {
	tokenDetails, err := srv.DBHelper.GetTokens()
	if err != nil {
		utils.EncodeJSONBody(resp, http.StatusInternalServerError, map[string]interface{}{
			"error": "unable to gte tokens ",
		})
		return
	}

	utils.EncodeJSON200Body(resp, map[string]interface{}{
		"tokenDetails": tokenDetails,
	})
}
