package main

import (
	"log"
	"net/http"

	"veracode.com/mypng/internal/api"
	"veracode.com/mypng/internal/web"
)

func main() {

	logger := log.Default()

	db, err := api.GetDB()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	apiSvc := api.Build(db)

	mux := http.NewServeMux()
	apiSvc.AddRoutes(mux)
	web.AddRoutes(mux)

	logger.Println("Serving api on http://localhost:8888")
	http.ListenAndServe(":8888", mux)
}
