package api

import (
	"encoding/json"
	"net/http"
)

func LoginRoutes(userSvc *UserSvc) *http.ServeMux {

	loginMux := http.NewServeMux()

	loginMux.HandleFunc("POST /login", func(rw http.ResponseWriter, req *http.Request) {

		body := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}

		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		isAuthenticated := userSvc.Authenticate(body.Username, body.Password)

		var user *User
		if isAuthenticated {
			user = userSvc.Authorize(body.Username)
		} else {
			userSvc.logger.Printf("failed to authenticate user %s", body.Username)
			user = nil
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
