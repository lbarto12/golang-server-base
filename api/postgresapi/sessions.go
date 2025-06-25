package postgresapi

import (
	"errors"
	"golang-server-base/api/postgresapi/models"
	"golang-server-base/api/postgresapi/postgresutils"
	"golang-server-base/api/webtokensapi"
	"time"
)

func SignUp(request models.Account) error {
	if request.UserName == "" {
		return errors.New("invalid sign-up request: missing parameter 'UserName'")
	}
	if request.Email == "" {
		return errors.New("invalid sign-up request: missing parameter 'Email'")
	}
	if request.Password == "" {
		return errors.New("invalid sign-up request: missing parameter 'Password'")
	}

	db, err := Database()
	if err != nil {
		return err
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// Check if account exists, fail if so
	var accountExists models.Exists
	err = tx.Get(&accountExists, "SELECT count(1) FROM accounts WHERE user_email = $1;", request.Email)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if accountExists.Found() {
		_ = tx.Rollback()
		return errors.New("account already exists")
	}

	// Hash password
	request.Password, err = postgresutils.HashPassword(request.Password)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// (encrypted) user_acc -> (Root) Only what is needed to manage a session -> UserProfile (user_Profile_id) //login/sign
	  // -> JWT { user_prof_id }
	// (encrypted) user_priv -> (user_priv_id) Bank Details, Payment token etc. (root_id FK) //payment
	// user_prof -> (user_prof_id) First Name, Last Name, Address //anytime

  // todo: Case insensitive matching
	// todo: remove user_name
	// Insert row
	query, err := tx.NamedQuery("INSERT INTO accounts (user_name, user_email, password_hash) VALUES (:user_name, :user_email, :password_hash)", request)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	err = query.Close()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func SignIn(request models.Account, jwtSessionToken string) (string, error) {
	// Try with token first if fields are null, do not allow sign in with token if any credentials are provided alongside it
	if jwtSessionToken != "" && request.Email == "" && request.Password == "" {
		token, err := webtokensapi.VerifyToken(jwtSessionToken)
		if err == nil && token.Valid {
			return jwtSessionToken, nil
		}
	}

	// Then with credentials
	if request.Email == "" {
		return "", errors.New("invalid request: missing parameter 'Email'")
	}

	if request.Password == "" {
		return "", errors.New("invalid request: missing parameter 'Password'")
	}

	db, err := Database()
	if err != nil {
		return "", err
	}

	tx, err := db.Beginx()
	if err != nil {
		return "", err
	}

	var account models.Account
	err = tx.Get(&account, "SELECT * FROM accounts WHERE user_email = $1;", request.Email)
	if err != nil {
		_ = tx.Rollback()
		if account.Email == "" {
			return "", errors.New("user not found")
		}
		return "", err
	}

	err = postgresutils.PasswordEqualsHash(request.Password, account.Password)
	if err != nil {
		_ = tx.Rollback()
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return "", err
	}

	return webtokensapi.GenerateJWT(request.Email, time.Now().Add(24*time.Hour))
}
