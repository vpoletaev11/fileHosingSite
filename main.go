package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vpoletaev11/fileHostingSite/cookie"
	"github.com/vpoletaev11/fileHostingSite/pages/index"
	"github.com/vpoletaev11/fileHostingSite/pages/login"
	"github.com/vpoletaev11/fileHostingSite/pages/logout"
	"github.com/vpoletaev11/fileHostingSite/pages/registration"
	"github.com/vpoletaev11/fileHostingSite/pages/upload"
)

func main() {
	// connecting to mySQL database
	db, err := sql.Open("mysql", "perdator:@tcp(localhost:3306)/fileHostingSite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Successfully connected to MySql database")

	// creating file server handler for assets
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// registration page handler
	http.HandleFunc("/register", registration.Page(db))

	// login page handler
	http.HandleFunc("/login", login.Page(db))

	// index page handler
	http.HandleFunc("/", cookie.AuthWrapper(index.Page(db), db))

	// logout page handler
	http.HandleFunc("/logout", logout.Page(db))

	// upload file page handler
	http.HandleFunc("/upload", cookie.AuthWrapper(upload.Page(db), db))

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
