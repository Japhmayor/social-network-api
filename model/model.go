package model

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email"`
	First string `json:"firstname"`
	Last  string `json:"lastname"`
}

func (u *User) GetUser(db *sql.DB) error {
	err := db.QueryRow("SELECT * FROM users WHERE user_id=?", u.ID).Scan(&u.Name, &u.Email, &u.First, &u.Last)
	return err
}

func (u *User) UpdateUser(db *sql.DB) error {
	_, err := db.Exec("UPDATE users SET user_nickname=?, user_email=?, user_first=?, user_last=? WHERE user_id=?", u.Name, u.Email, u.First, u.Last)
	return err
}

func (u *User) DeleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users WHERE user_id=?", u.ID)
	return err
}

func (u *User) CreateUser(db *sql.DB) error {
	// err := db.QueryRow("INSERT INTO user(name, email, first, last) VALUES(?, ?, ?, ?) RETURNING id", u.Name, u.Email, u.First, u.Last).Scan(&u.ID)
	sql := "INSERT INTO users set user_nickname='" + u.Name + "', user_first='" + u.First + "', user_last='" + u.Last + "', user_email='" + u.Email + "'"
	result, err := db.Exec(sql)
	if err != nil {
		return err
	}
	fmt.Print(result)
	return nil
}

func GetUsers(db *sql.DB) []User {
	rows, err := db.Query("select * from users LIMIT 10")
	if err != nil {
		return nil
	}
	defer rows.Close()
	users := []User{}
	for rows.Next() {
		u := User{}
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.First, &u.Last)
		users = append(users, u)

	}
	return users
}
