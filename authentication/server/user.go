package server

import (
	"Assignment/models"
	"Assignment/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// VerifyUser - api to validate user using invite code generated by admin
// taking invite code as query param
// storing invite code in in-memory cache
// fetching code from in-memory cache
// if present then user is logged in else need try again
func (srv *Server) VerifyUser(resp http.ResponseWriter, req *http.Request) {
	ipAddress := req.Header.Get("X-REAL-IP")
	token := req.URL.Query().Get("token")
	if token == "" {
		utils.EncodeJSONBody(resp, http.StatusBadRequest, map[string]interface{}{
			"error": errors.New("invite token cannot be empty"),
		})
		return
	}

	if ipAddress == "" {
		ipAddress = req.RemoteAddr
	}

	expirationTime := time.Now().Add(2 * time.Minute)
	var tokenDetails models.TokenDetail
	byteData, err := srv.CacheProvider.Get(ipAddress)
	if err == nil {
		err = json.Unmarshal([]byte(fmt.Sprintf("%v", byteData)), &tokenDetails)
		if err != nil {
			utils.EncodeJSONBody(resp, http.StatusInternalServerError, map[string]interface{}{
				"error": err,
			})
			return
		}

		if tokenDetails.Count > 5 {
			utils.EncodeJSONBody(resp, http.StatusForbidden, map[string]interface{}{
				"message": "Too Many Attempts! Wait for few minutes",
			})
			return
		}

		expirationTime = tokenDetails.ExpirationTime

	}

	err = srv.CacheProvider.Set(ipAddress, &models.TokenDetail{Token: token, Count: tokenDetails.Count + 1, ExpirationTime: expirationTime})
	if err != nil {
		logrus.Errorf("unable to cache %v", err)
	}

	fmt.Print(srv.CacheProvider.Get(ipAddress))
	_, err = srv.CacheProvider.Get(token)
	if err == nil {
		utils.EncodeJSON200Body(resp, map[string]interface{}{
			"message": "successful login",
		})
		return
	}

	////using database
	//isValid, err := srv.DBHelper.VerifyToken(token)
	//if err != nil {
	//	utils.EncodeJSONBody(resp, http.StatusInternalServerError, map[string]interface{}{
	//		"error": err,
	//	})
	//	return
	//}
	//
	//if isValid {
	//	utils.EncodeJSON200Body(resp, map[string]interface{}{
	//		"message": "successful login",
	//	})
	//	return
	//}

	utils.EncodeJSONBody(resp, http.StatusUnauthorized, map[string]interface{}{
		"message": "please retry",
	})
}
