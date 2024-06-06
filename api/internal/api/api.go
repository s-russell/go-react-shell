package api

import (
	"github.com/jmoiron/sqlx"
	"net/http"
	"veracode.com/mypng/internal/api/routes"
	"veracode.com/mypng/internal/api/services"
)

type API struct {
	UserSvc *services.UserSvc
}

func Build(db *sqlx.DB) API {
	userSvc := services.NewUserSvc(db)
	return API{&userSvc}
}

func (api *API) AddRoutes(mux *http.ServeMux) {

	mux.Handle("/api/login/*", http.StripPrefix("/api/login", routes.Login(api.UserSvc)))
	mux.Handle("/api/user/*", http.StripPrefix("/api/user", routes.User(api.UserSvc)))

}
