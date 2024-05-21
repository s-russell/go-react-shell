package api

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Name      string   `json:"name"`
	Roles     []string `json:"roles"`
}

type UserSvc struct {
	logger *log.Logger
	db     *sqlx.DB
}

func NewUserSvc(db *sqlx.DB) UserSvc {
	logger := log.New(os.Stdout, "UserSvc: ", log.LstdFlags|log.Lshortfile)
	return UserSvc{logger, db}
}

func (userSvc *UserSvc) Create(user *User, password string) (int64, error) {
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return -1, err
	}

	query := `
	insert into user(firstName, lastName, username, password)
	values ($1, $2, $3, $4)
	`

	result, err := userSvc.db.Exec(query, user.FirstName, user.LastName, user.Name, passwordHashed)
	if err != nil {
		return -1, err
	}

	insertedId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	userSvc.logger.Printf("createdd user %s", user.Name)
	return insertedId, nil
}

func (userSvc *UserSvc) Authenticate(username string, password string) bool {

	query := `
	select password from user 
		where username = $1
	`

	rows, err := userSvc.db.Query(query, username)
	if err != nil {
		userSvc.logger.Printf("error authenticating user %s:\n%s\n", username, err)
		return false
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			userSvc.logger.Printf("error cleaning up when authenticating user: %s\n", err)
		}
	}(rows)

	if rows.Next() {
		var lastPassword string
		if err := rows.Scan(&lastPassword); err != nil {
			userSvc.logger.Printf("error authenticating user %s:\n%s\n", username, err)
			return false
		}
		err = bcrypt.CompareHashAndPassword([]byte(lastPassword), []byte(password))
		return err == nil
	}

	return false
}

func (userSvc *UserSvc) Authorize(username string) *User {

	query := `
             select u.first_name, u.last_name, r.name "role_name"
             from user u,
                  user_role ur,
                  roles r
             where u.username = $1
               and u.id = ur.user_id
               and r.id = ur.role_id
             `
	user := User{}

	rows, err := userSvc.db.Queryx(query, username)
	if err != nil {
		userSvc.logger.Printf("error retrieving roles for user %s:\n%s\n", username, err)
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			userSvc.logger.Printf("error cleaning up when retrieving roles for user  %s:\n%s\n", username, err)
		}
	}(rows)

	var firstName, lastName, role string
	for rows.Next() {
		err = rows.Scan(&firstName, &lastName, &role)
		if err != nil {
			userSvc.logger.Printf("error retrieving roles for user %s:\n%s\n", username, err)
		}
		user.FirstName = firstName
		user.LastName = lastName
		user.Roles = append(user.Roles, role)
	}

	user.Name = username

	return &user

}
