package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	err := initializeTransactionLog()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.HandleFunc("/", healthcheckHandler)
	r.HandleFunc("/v1/key/{key}", keyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/key/{key}", keyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/key/{key}", keyValueDeleteHandler).Methods("DELETE")

	fmt.Println("Listening on port 8000")

	log.Fatal(http.ListenAndServe(":8000", r))
}
