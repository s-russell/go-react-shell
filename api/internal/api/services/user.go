package services

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"database/sql"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	sessionLength = 1 * time.Hour
)

type UserClaims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

func (claims *UserClaims) hasRole(candidate string) bool {
	for _, role := range claims.Roles {
		if candidate == role {
			return true
		}
	}
	return false
}

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
	token := jwt.NewWithClaims(jwt.SigningMethodES256, &UserClaims{
		jwt.RegisteredClaims{
			Issuer:    "steve",
			Subject:   user.Name,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(sessionLength)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		user.Roles,
	})

	return token.SignedString(userSvc.secret)
}

type httpLink func(rw http.ResponseWriter, r *http.Request)

func (userSvc *UserSvc) AuthorizeAnyRoleMiddleware(roles ...string) func(httpLink) http.Handler {
	return func(nextHandlerFunc httpLink) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			//Authenticate
			claims, ok := validateToken(r, &userSvc.secret.PublicKey)
			if !ok {
				http.Error(rw, "Forbidden", http.StatusForbidden)
				return
			}

			//Authorize
			isAuthorized := false
			for _, role := range roles {
				isAuthorized = isAuthorized || claims.hasRole(role)
			}

			if !isAuthorized {
				http.Error(rw, "Unauthorized", http.StatusUnauthorized)
				return
			}

			http.HandlerFunc(nextHandlerFunc).ServeHTTP(rw, r)
		})
	}
}

func validateToken(r *http.Request, publicKey *ecdsa.PublicKey) (*UserClaims, bool) {

	tokenCookie, err := r.Cookie("session_token")
	if err != nil || tokenCookie == nil {
		return nil, false
	}

	token, err := jwt.Parse(tokenCookie.String(), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(UserClaims); !ok {
		return nil, false
	} else {
		return &claims, true
	}
}
