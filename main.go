package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vpoletaev11/fileHostingSite/cookie"
	"github.com/vpoletaev11/fileHostingSite/pages/admin"
	"github.com/vpoletaev11/fileHostingSite/pages/categories"
	"github.com/vpoletaev11/fileHostingSite/pages/cookiescleaner"
	"github.com/vpoletaev11/fileHostingSite/pages/download"
	"github.com/vpoletaev11/fileHostingSite/pages/index"
	"github.com/vpoletaev11/fileHostingSite/pages/login"
	"github.com/vpoletaev11/fileHostingSite/pages/logout"
	"github.com/vpoletaev11/fileHostingSite/pages/popular"
	"github.com/vpoletaev11/fileHostingSite/pages/registration"
	"github.com/vpoletaev11/fileHostingSite/pages/upload"
	"github.com/vpoletaev11/fileHostingSite/pages/users"
)

func main() {
	// connecting to mySQL database
	db, err := sql.Open("mysql", "root:@tcp(mysql:3306)/fileHostingSite?parseTime=true") // ?parseTime=true asks the driver to scan DATE and DATETIME automatically to time.Time
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Successfully connected to MySql database")

	// creating file server handler for assets
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// creating file server handler for files
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("files"))))

	http.HandleFunc("/registration", registration.Page(db))
	http.HandleFunc("/login", login.Page(db))
	http.HandleFunc("/", cookie.AuthWrapper(index.Page, db))
	http.HandleFunc("/logout", logout.Page(db))
	http.HandleFunc("/upload", cookie.AuthWrapper(upload.Page, db))
	http.HandleFunc("/categories/", cookie.AuthWrapper(categories.Page, db))
	http.HandleFunc("/download", cookie.AuthWrapper(download.Page, db))
	http.HandleFunc("/popular", cookie.AuthWrapper(popular.Page, db))
	http.HandleFunc("/users", cookie.AuthWrapper(users.Page, db))
	http.HandleFunc("/admin", cookie.AdminAuthWrapper(admin.Page, db))
	http.HandleFunc("/cookiescleaner", cookiescleaner.Page(db))

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
