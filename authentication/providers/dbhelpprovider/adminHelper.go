package dbhelpprovider

import (
	"Assignment/models"
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
)

func (dh *DBHelper) GetTokens() ([]models.TokenDetails, error) {

	tokenList := make([]models.TokenDetails, 0)
	// language=sql
	SQL := `SELECT invite_token,
                   COALESCE(archived_at < now(), false) as is_archived
            FROM token_info
            order by created_at desc`

	err := dh.DB.Select(&tokenList, SQL)
	if err != nil && err != sql.ErrNoRows {
		logrus.Errorf("GetTokens: error in getting token details: %v", err)
		return tokenList, err
	}

	return tokenList, nil
}

func (dh *DBHelper) LogInUserUsingEmail(loginReq models.LoginRequest) (userID string, message string, err error) {
	// language=SQL
	SQL := `
		SELECT 
		    id,   
			password
		FROM
			users
		WHERE
			email = $1
			AND archived_at IS NULL 
	`

	var user = struct {
		ID       string `db:"id"`
		Password string `db:"password"`
	}{}

	if err = dh.DB.Get(&user, SQL, loginReq.Email); err != nil && err != sql.ErrNoRows {
		logrus.Errorf("LogInUserUsingEmail: error while getting user %v", err)
		return userID, "error getting user", err
	}

	if user.Password != loginReq.Password {
		return userID, "Password Not Correct", errors.New("password not matched")
	}

	return user.ID, "", nil
}

func (dh *DBHelper) AddToken(inviteTokenDetail models.TokenDetail) error {
	//	language=sql
	SQL := `INSERT INTO token_info(invite_token, archived_at)
			VALUES ($1,$2)`

	_, err := dh.DB.Exec(SQL, inviteTokenDetail.Token, inviteTokenDetail.ExpirationTime)
	if err != nil {
		logrus.Errorf("AddToken: error in inserting token details: %v", err)
		return err
	}
	return nil
}

func (dh *DBHelper) DeleteToken(tokenID string) error {
	// language=sql
	SQL := `UPDATE token_info set archived_at = now() where invite_token = $1`

	_, err := dh.DB.Exec(SQL, tokenID)
	if err != nil {
		logrus.Errorf("DeleteToken: error in deleting token details: %v", err)
		return err
	}
	return nil
}
