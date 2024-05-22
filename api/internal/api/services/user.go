package services

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"database/sql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
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
	Logger *log.Logger
	db     *sqlx.DB
	secret *ecdsa.PrivateKey // used to sign JWTs--restarting server logs everyone out
}

func NewUserSvc(db *sqlx.DB) UserSvc {
	logger := log.New(os.Stdout, "UserSvc: ", log.LstdFlags|log.Lshortfile)

	secret, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logger.Fatalf("Could not generate ECDSA key for request authentication %v", err)
	}

	return UserSvc{
		logger,
		db,
		secret,
	}
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
	userSvc.Logger.Printf("createdd user %s", user.Name)
	return insertedId, nil
}

func (userSvc *UserSvc) Authenticate(username string, password string) bool {

	query := `
	select password from user 
		where username = $1
	`

	rows, err := userSvc.db.Query(query, username)
	if err != nil {
		userSvc.Logger.Printf("error authenticating user %s:\n%s\n", username, err)
		return false
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			userSvc.Logger.Printf("error cleaning up when authenticating user: %s\n", err)
		}
	}(rows)

	if rows.Next() {
		var lastPassword string
		if err := rows.Scan(&lastPassword); err != nil {
			userSvc.Logger.Printf("error authenticating user %s:\n%s\n", username, err)
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
		userSvc.Logger.Printf("error retrieving roles for user %s:\n%s\n", username, err)
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			userSvc.Logger.Printf("error cleaning up when retrieving roles for user  %s:\n%s\n", username, err)
		}
	}(rows)

	var firstName, lastName, role string
	for rows.Next() {
		err = rows.Scan(&firstName, &lastName, &role)
		if err != nil {
			userSvc.Logger.Printf("error retrieving roles for user %s:\n%s\n", username, err)
		}
		user.FirstName = firstName
		user.LastName = lastName
		user.Roles = append(user.Roles, role)
	}

	user.Name = username

	return &user

}

func (userSvc *UserSvc) BuildJWT(user *User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub":       user.Name,
		"roles":     user.Roles,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
	})

	return token.SignedString(userSvc.secret)
}

type httpLink func(rw http.ResponseWriter, r *http.Request)

func (userSvc *UserSvc) AuthorizeRolesMiddleware(roles ...string) func(httpLink) http.Handler {
	return func(nextHandlerFunc httpLink) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if !requestIsAuthenticated(r) {
				http.Error(rw, "Forbidden", http.StatusForbidden)
			}

			if !requestIsAuthorized(userSvc, roles...) {
				http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			}

			http.HandlerFunc(nextHandlerFunc).ServeHTTP(rw, r)
		})
	}
}

func requestIsAuthenticated(r *http.Request) bool {

	authHeader, err := r.Cookie("session_token")
	if err != nil {
		return false
	}

	// TODO: actually validate the token
	authHeader = authHeader

	return true
}

func requestIsAuthorized(svc *UserSvc, roles ...string) bool {
	return true
}
