package api

import (
	"github.com/jmoiron/sqlx"
	"net/http"
)

type API struct {
	UserSvc *UserSvc
}

func Build(db *sqlx.DB) API {
	userSvc := NewUserSvc(db)
	return API{&userSvc}
}

func (api *API) AddRoutes(mux *http.ServeMux) {

	mux.Handle("POST /api/", http.StripPrefix("/api", LoginRoutes(api.UserSvc)))

}
