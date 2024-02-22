package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Data struct {
	Id       string `json:"id"`
	Lastname string `json:"lastname"`
}

var (
	Db *sql.DB
)

func main() {
	r := mux.NewRouter()
	fmt.Println("here")
	d, err := sql.Open("postgres", URLDatabase())
	if err != nil {
		log.Fatal(err)
	}
	Db = d
	_, err = Db.Exec("CREATE TABLE IF NOT EXISTS testt (ID integer, lastname varchar(255))")
	if err != nil {
		fmt.Println("here1")
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "<html><head><title>My Golang Web App</title></head><body><h1>Hello, Golang Web App!</h1></body></html>")
	}).Methods("GET")
	r.HandleFunc("/getdata", GetData).Methods("GET")
	http.ListenAndServe(port, r)
}
func GetData(w http.ResponseWriter, r *http.Request) {
	rows, err := Db.Query("SELECT * FROM testT")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var data []Data
	for rows.Next() {
		var id, lname string
		err := rows.Scan(&id, &lname)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Data{Id: id, Lastname: lname})
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
func URLDatabase() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}
