package middlewareprovider

import (
	"Assignment/models"
	"Assignment/utils"
	"context"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

const (
	authorization = "Authorization"
	bearerScheme  = "bearer"
	space         = " "
	maxAge        = 300
)

func corsOptions() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Token", "importDate", "X-Client-Version", "Cache-Control", "Pragma", "x-session-token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           maxAge, // Maximum value not ignored by any of major browsers
	})
}

func authMiddleware(db *sqlx.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			mySigningKey := []byte(os.Getenv("JWT_SECRET_KEY"))
			jwtKey := req.Header.Get("Authorization")
			fmt.Printf("got token: %s\n", jwtKey)
			decryptedToken, err := jwt.Parse(jwtKey, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return mySigningKey, nil
			})
			if claims, ok := decryptedToken.Claims.(jwt.MapClaims); ok && decryptedToken.Valid {
				userID := claims["userID"].(string)

				var userContext *models.UserContext

				userContext, err = getUserByID(db, userID)
				if err != nil {
					logrus.Errorf("authMiddleware: error in getting user info: %v", err)
					utils.EncodeJSONBody(resp, http.StatusUnauthorized, "user not authorised")
					return
				}

				req = req.WithContext(context.WithValue(req.Context(), models.UserContextKey, userContext))
				next.ServeHTTP(resp, req)
			} else {
				fmt.Println(err)
				utils.EncodeJSONBody(resp, http.StatusUnauthorized, "user not authorised")
				return
			}
		})
	}
}

func checkForEmptyContext() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			_, ok := req.Context().Value(models.UserContextKey).(*models.UserContext)
			if !ok {
				utils.EncodeJSONBody(resp, http.StatusUnauthorized, "user not authorized")
				return
			}
			next.ServeHTTP(resp, req)
		})
	}
}

func getUserByID(db *sqlx.DB, userID string) (*models.UserContext, error) {
	SQL := `SELECT 
				id,
				name,
				email
			 FROM 
				users 
			 WHERE 
				id = $1 AND 
			     archived_at is null  
				`

	var userCtx models.UserContext

	err := db.Get(&userCtx, SQL, userID)
	if err != nil {
		logrus.Errorf("getUserByID: error getting user context %v", err)
		return nil, err
	}

	return &userCtx, nil
}
