package models

import (
	"database/sql"

	"github.com/learninNdi/go-session-login-register/config"
	"github.com/learninNdi/go-session-login-register/entities"
)

type UserModel struct {
	db *sql.DB
}

func NewUserModel() *UserModel {
	conn, err := config.DBConn()

	if err != nil {
		panic(err)
	}

	return &UserModel{
		db: conn,
	}
}

func (u UserModel) Where(user *entities.User, fieldName, fieldValue string) error {
	row, err := u.db.Query("select id, username, fullname, email, password from users where "+fieldName+" = ? limit 1", fieldValue)

	if err != nil {
		return err
	}

	defer row.Close()

	for row.Next() {
		row.Scan(&user.ID, &user.Username, &user.FullName, &user.Email, &user.Password)
	}

	return nil
}
