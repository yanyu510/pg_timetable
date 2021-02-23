package api

import (
	"net/http"

	"github.com/cybertec-postgresql/pg_timetable/internal/pgengine"
)

// InitWebServer for web api
func InitWebServer(retrive chan int, port string) {

	// Handler request of retrive
	http.HandleFunc("/retrive", func(w http.ResponseWriter, r *http.Request) {
		retrive <- 1
	})
	pgengine.Log("LOG", "Webserver listening in ", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		pgengine.Log("ERROR", "Webserver failed.%s", err.Error())
	}
}
