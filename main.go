package main

import (
	// "database/sql"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var database *sql.DB

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	First string `json:"first"`
	Last  string `json:"last"`
}

type Users struct {
	Users []User `json:"users"`
}

func UsersCreate(w http.ResponseWriter, r *http.Request) {
	NewUser := User{}
	NewUser.Name = r.FormValue("user")
	NewUser.Email = r.FormValue("email")
	NewUser.First = r.FormValue("first")
	NewUser.Last = r.FormValue("last")

	sql := "INSERT INTO users set user_nickname='" + NewUser.Name +
		"', user_first='" + NewUser.First + "', user_last='" +
		NewUser.Last + "', user_email='" + NewUser.Email + "'"
	q, err := database.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(q)
}

func UsersRetrieve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("pragma", "no-cache")
	rows, _ := database.Query("SELECT * FROM users LIMIT 10")
	Response := Users{}
	for rows.Next() {
		user := User{}
		rows.Scan(&user.ID, &user.Name, &user.First, &user.Last, &user.Email)
		Response.Users = append(Response.Users, user)
	}
	output, _ := json.Marshal(Response)
	fmt.Fprintf(w, string(output))
}

func main() {
	db, err := sql.Open("mysql", "root:bontusfavor1994?@tcp(127.0.0.1:3306)/social_network")
	database = db
	if err != nil {
		fmt.Println("cannot connect to database social_network")
		log.Fatal("error", err)
	}
	r := mux.NewRouter()
	// routes
	r.HandleFunc("/api/users", UsersRetrieve).Methods("GET")
	r.HandleFunc("/api/users", UsersCreate).Methods("POST")
	http.Handle("/", r)
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Print("Error", err)
	}
}
