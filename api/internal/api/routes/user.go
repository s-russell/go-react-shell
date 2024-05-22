package routes

import (
	"net/http"
	"veracode.com/mypng/internal/api/services"
)

func User(userSvc *services.UserSvc) *http.ServeMux {

	userMux := http.NewServeMux()

	devsOnly := userSvc.AuthorizeRolesMiddleware("developers")

	userMux.Handle("get /user/devsonly", devsOnly(func(rw http.ResponseWriter, req *http.Request) {
		return
	}))

	return userMux
}
