package routes

import (
	"jerubaal.com/goreactor/internal/api/services"
	"net/http"
)

func User(userSvc *services.UserSvc) *http.ServeMux {

	userMux := http.NewServeMux()

	devsOnly := userSvc.AuthorizeAnyRoleMiddleware("developers")

	userMux.Handle("GET /devsonly", devsOnly(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusAccepted)
	}))

	return userMux
}
