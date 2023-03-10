package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func handleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		email := r.PostFormValue("email")
		password := r.PostFormValue("password")

		rows, err := db.Query("SELECT password FROM public.User WHERE email=$1", email)

		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var password_matched bool

		for rows.Next() {
			var _password string

			err = rows.Scan(&_password)

			if err != nil {
				panic(err)
			}

			password_matched = password == _password
		}

		if password_matched {
			http.SetCookie(w, &http.Cookie{
				Name:    "session",
				Value:   email,
				Expires: time.Now().Add(5 * time.Minute),
			})
			w.Header().Set("Goto", "/menu.html")
			w.WriteHeader(http.StatusSeeOther)
		} else {
			w.Header().Set("Content-Type", "application/json")
			data := make(map[string]string)
			data["reason"] = "Invalid Credentials"
			json.NewEncoder(w).Encode(data)
		}
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		email := r.PostFormValue("email")
		password := r.PostFormValue("password")
		confirm := r.PostFormValue("password-confirm")

		if password != confirm {
			w.Header().Set("Content-Type", "application/json")
			data := make(map[string]string)
			data["reason"] = "Password don't match"
			json.NewEncoder(w).Encode(data)
		} else {
			_, err := db.Exec(`INSERT INTO public.User VALUES ($1, $2)`, email, password)

			if err != nil {
				panic(err)
			}

			w.Header().Set("Goto", "/login.html")
			w.WriteHeader(http.StatusSeeOther)
		}
	}
}

func handlePurchaseItem(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		cookie, errc := r.Cookie("session")

		if errc != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user_email := cookie.Value
		index := r.PostFormValue("index")
		amount := r.PostFormValue("amount")

		_, err := db.Exec(`CALL purchase_item($1, $2, $3)`, user_email, index, amount)

		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func handlePurchaseItemCount(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		cookie, errc := r.Cookie("session")

		if errc != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		email := cookie.Value

		rows, err := db.Query("SELECT index, count FROM public.Item WHERE user_email=$1", email)

		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var columns [][]int

		for rows.Next() {
			cols := make([]int, 2)

			if err := rows.Scan(&cols[0], &cols[1]); err != nil {
				panic(err)
			}

			columns = append(columns, cols)
		}

		json.NewEncoder(w).Encode(columns)
	}
}

func main() {
	var (
		host     = os.Getenv("POSTGRES_HOST")
		port     = os.Getenv("POSTGRES_PORT")
		user     = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbname   = os.Getenv("POSTGRES_DB")
	)

	connection_string := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error

	db, err = sql.Open("postgres", connection_string)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/create-account", handleRegister)
	http.HandleFunc("/validate-account", handleLogin)
	http.HandleFunc("/purchase-item", handlePurchaseItem)
	http.HandleFunc("/purchase-item-count", handlePurchaseItemCount)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
	fmt.Println("Listening on port: " + os.Getenv("PORT"))
}
