package routes

import (
	"encoding/json"
	"net/http"
	"veracode.com/mypng/internal/api/services"
)

func Login(userSvc *services.UserSvc) *http.ServeMux {

	loginMux := http.NewServeMux()

	loginMux.HandleFunc("POST /", func(rw http.ResponseWriter, req *http.Request) {

		body := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}

		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		isAuthenticated := userSvc.Authenticate(body.Username, body.Password)

		var user *services.User
		if isAuthenticated {
			user = userSvc.Authorize(body.Username)
		} else {
			userSvc.Logger.Printf("failed to authenticate user %s", body.Username)
			user = nil
		}

		if user != nil {
			if token, err := userSvc.BuildJWT(user); err != nil {
				userSvc.Logger.Printf("failed to build JWT for user %s: %v", user.Name, err)
			} else {
				http.SetCookie(rw, &http.Cookie{
					Name:     "session_token",
					Value:    token,
					Path:     "/",  // Available to the entire domain
					HttpOnly: true, // Makes the cookie inaccessible to client-side scripts
					MaxAge:   3600, // Expires after 1 hour
				})
			}

		}

		rw.Header().Set("Content-Type", "application/json")

		jsonData, err := json.Marshal(user)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Write(jsonData)
	})

	return loginMux

}
