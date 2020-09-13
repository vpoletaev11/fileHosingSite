package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/vpoletaev11/fileHostingSite/pages/categories"
	"github.com/vpoletaev11/fileHostingSite/pages/download"
	"github.com/vpoletaev11/fileHostingSite/pages/index"
	"github.com/vpoletaev11/fileHostingSite/pages/login"
	"github.com/vpoletaev11/fileHostingSite/pages/logout"
	"github.com/vpoletaev11/fileHostingSite/pages/popular"
	"github.com/vpoletaev11/fileHostingSite/pages/registration"
	"github.com/vpoletaev11/fileHostingSite/pages/upload"
	"github.com/vpoletaev11/fileHostingSite/pages/users"
	"github.com/vpoletaev11/fileHostingSite/session"
)

const (
	mySQLAddr = "user:@tcp(localhost:3306)"
	redisAddr = "localhost:6379"
)

func main() {
	dep := connectToDBs()

	// creating file server handler for assets
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// creating file server handler for files
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("files"))))

	http.HandleFunc("/registration", registration.Page(dep.Db))
	http.HandleFunc("/login", login.Page(dep))
	http.HandleFunc("/", session.AuthWrapper(index.Page, dep))
	http.HandleFunc("/logout", logout.Page(dep))
	http.HandleFunc("/upload", session.AuthWrapper(upload.Page, dep))
	http.HandleFunc("/categories/", session.AuthWrapper(categories.Page, dep))
	http.HandleFunc("/download", session.AuthWrapper(download.Page, dep))
	http.HandleFunc("/popular", session.AuthWrapper(popular.Page, dep))
	http.HandleFunc("/users", session.AuthWrapper(users.Page, dep))

	fmt.Println("Starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

func connectToDBs() session.Dependency {
	// connecting to mySQL database
	db, err := sql.Open("mysql", mySQLAddr+"/fileHostingSite?parseTime=true") // ?parseTime=true asks the driver to scan DATE and DATETIME automatically to time.Time
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to MySql database")

	// connecting to Redis
	redisConn, err := redis.Dial("tcp", redisAddr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to Redis")

	// Connections to databases will be closed after program exit
	return session.Dependency{Db: db, Redis: redisConn}
}
