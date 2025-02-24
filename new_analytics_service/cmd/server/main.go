package main

import (
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	err := http.ListenAndServe(":8002", router)
	if err != nil {
		log.Fatal("HTTP server failed.")
		return
	}
}
