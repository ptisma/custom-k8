package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/version", getVersionHandler).Methods("GET")

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	fmt.Println("Server is listening on :8080...")
	log.Fatal(server.ListenAndServe())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := &Response{Payload: "Hello world"}
	jData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func getVersionHandler(w http.ResponseWriter, r *http.Request) {
	ver := os.Getenv("VERSION")
	msg := fmt.Sprintf("Version:%s", ver)
	res := Response{Payload: msg}
	jsonData, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

type Response struct {
	Payload string `json:"payload"`
}
