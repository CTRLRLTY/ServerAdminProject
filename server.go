package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handleRegister(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// email := r.PostFormValue("email")
		password := r.PostFormValue("password")
		confirm := r.PostFormValue("password-confirm")

		if password != confirm {
			w.Header().Set("Content-Type", "application/json")
			data := make(map[string]string)
			data["reason"] = "Password don't match"
			json.NewEncoder(w).Encode(data)
		} else {
			w.Header().Set("Goto", "/register.js")
			w.WriteHeader(http.StatusSeeOther)
		}
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/create-account", handleRegister)

	fmt.Println("Starting!")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
