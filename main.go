package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type Book struct {
	ID		int		`json:"id"`
	Title	string	`json:"title"`
	Author	string	`json:"author"`
	Year	string	`json:"year"`
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func dbConn() (db *gorm.DB) {
	dbhost     := os.Getenv("DB_HOST")
	dbport     := os.Getenv("DB_PORT")
	dbuser     := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname     := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpassword, dbname)

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	db := dbConn()
	db.AutoMigrate(&Book{})
	db.Close()
	db.LogMode(true)

	r := mux.NewRouter()

	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	//r.HandleFunc("/books", addBook).Methods("POST")
	//r.HandleFunc("/books", updateBook).Methods("PUT")
	//r.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	log.Printf("Starting server on %s", port)

	http.ListenAndServe(":" + port, r)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	books := []Book{}
	db.Find(&books)
	db.Close()

	json, err := json.Marshal(books)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	books := []Book{}
	params := mux.Vars(r)
	id := params["id"]

	db.First(&books, id)
	db.Close()

	json, err := json.Marshal(books)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func addBook(w http.ResponseWriter, r *http.Request) {

}

//func updateBook(w http.ResponseWriter, r *http.Request) {
//
//}
//
//func removeBook(w http.ResponseWriter, r *http.Request) {
//
//}
