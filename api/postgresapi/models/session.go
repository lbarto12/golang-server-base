package models

type Account struct {
	Id       int    `db:"user_id"`
	UserName string `db:"user_name"`
	Email    string `db:"user_email"`
	Password string `db:"password_hash"`
}
