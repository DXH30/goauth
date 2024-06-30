package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "html/template"
    _ "github.com/mattn/go-sqlite3"
)

type User struct {
    ID int
    Username string
    Password string
}

var db *sql.DB

func initDB() {
    fmt.Println("Initialize Database...")
    var err error
    db, err = sql.Open("sqlite3", "./user.db")
    if err != nil {
        log.Fatal(err)
    }

    createTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );
    `

    _, err = db.Exec(createTable)
    if err != nil {
        log.Fatal(err)
    }
}

func registerUser(username, password string) error {
    stmt, err := db.Prepare("INSERT INTO users(username, password) VALUES(?, ?)")
    if err != nil {
        return err
    }

    _, err = stmt.Exec(username, password)
    if err != nil {
        return err
    }

    return nil
}

func authenticateUser(username, password string) bool {
    var user User
    err := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password)
    if err != nil {
        return false
    }
    return user.Password == password
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        tmpl, _ := template.ParseFiles("login.html")
        tmpl.Execute(w, nil)
        return
    }

    if r.Method == http.MethodPost {
        r.ParseForm()
        username := r.FormValue("username")
        password := r.FormValue("password")

        if authenticateUser(username, password) {
            fmt.Fprintf(w, "Login successful")
        } else {
            fmt.Fprintf(w, "Invalid username or password")
        }
    }
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        tmpl, _ := template.ParseFiles("register.html")
        tmpl.Execute(w, nil)
        return
    }

    if r.Method == http.MethodPost {
        r.ParseForm()
        username := r.FormValue("username")
        password := r.FormValue("password")

        err := registerUser(username, password)
        if err != nil {
            fmt.Fprintf(w, "Error: %v", err)
        } else {
            fmt.Fprintf(w, "Registration successful")
        }
    }
}

func main() {
    initDB()
    defer db.Close()

    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/register", registerHandler)

    fmt.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
