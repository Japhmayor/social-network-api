package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/japhmayor/social-media-api/model"
)

type App struct {
	Router *mux.Router
	Db     *sql.DB
}

// This initializes your database. It takes in the database user, password and the db name and returns an error.
func (a *App) Initialize(user, password, host, port, dbName string) {
	var err error
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, dbName)
	fmt.Println("Connecting to database", dbName, "...")
	a.Db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database connected")

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) Run(addr string) {
	fmt.Println("Starting on port", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// routes
func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uID, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user id")
		return
	}
	u := model.User{ID: uID}
	if err := u.GetUser(a.Db); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJson(w, http.StatusOK, u)
}

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("pragma", "no-cache")

	users := model.GetUsers(a.Db)

	respondWithJson(w, http.StatusOK, users)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var u model.User
	u.Name = r.FormValue("user")
	u.Email = r.FormValue("email")
	u.First = r.FormValue("first")
	u.Last = r.FormValue("last")

	if err := u.CreateUser(a.Db); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, u)

}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
	}
	u := model.User{ID: id}
	if err := u.DeleteUser(a.Db); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, u)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var u model.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	u.ID = id
	if err := u.UpdateUser(a.Db); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, u)

}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/api/user/{id:[0-9]+}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/api/user", a.createUser).Methods("POST")
	a.Router.HandleFunc("/api/user/[id:0-9]+", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/api/user/[id:0-9]+", a.deleteUser).Methods("DELETE")
}
