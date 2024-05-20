package api

import (
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

func (userAPI *UserSvc) Create(user *User, password string) (int64, error) {
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return -1, err
	}

	query := `
	insert into user(first_name, last_name, username, password)
	values ($1, $2, $3, $4)
	`

	result, err := userAPI.db.Exec(query, user.FirstName, user.LastName, user.Name, passwordHashed)
	if err != nil {
		return -1, err
	}

	insertedId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	userAPI.logger.Printf("createdd user %s", user.Name)
	return insertedId, nil
}

func (userAPI *UserSvc) Authenticate(username string, password string) bool {

	query := `
	select password from user 
		where username = $1
	`

	rows, err := userAPI.db.Query(query, username)
	if err != nil {
		userAPI.logger.Printf("error authenticating user %s:\n%s\n", username, err)
		return false
	}

	defer rows.Close()

	if rows.Next() {
		var lastPassword string
		if err := rows.Scan(&lastPassword); err != nil {
			userAPI.logger.Printf("error authenticating user %s:\n%s\n", username, err)
			return false
		}
		err = bcrypt.CompareHashAndPassword([]byte(lastPassword), []byte(password))
		return err == nil
	}

	return false
}

func (userSvc *UserSvc) Authorize(username string) *User {

	return &User{
		Name:  username,
		Roles: []string{"developer"},
	}
}
