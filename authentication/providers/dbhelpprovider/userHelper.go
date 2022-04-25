package dbhelpprovider

import (
	"Assignment/models"
	"github.com/sirupsen/logrus"
)

func (dh *DBHelper) VerifyToken(token string) (bool, error) {
	// language=sql
	SQL := `SELECT count(*) > 0
            FROM token_info 
            WHERE invite_token= $1`

	var isValid bool
	err := dh.DB.Get(&isValid, SQL, token)
	if err != nil {
		logrus.Errorf("verifyToken: enable toi veridy token %v", err)
		return isValid, err
	}

	return isValid, nil
}

func (dh *DBHelper) PopulateCache() ([]models.TokenDetail, error) {
	var tokenDetail []models.TokenDetail

	SQL := `SELECT invite_token , archived_at
	        FROM token_info
	        WHERE archived_at > now()`
	err := dh.DB.Select(&tokenDetail, SQL)
	if err != nil {
		logrus.Errorf("unable to cache tokens %v", err)
		return tokenDetail, err
	}

	return tokenDetail, nil
}
