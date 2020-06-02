package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type makanan struct {
	Namamakanan   string
	Jenishidangan string
	Hargamakanan  int
}

type view struct {
	Judul string
	Data  map[int]makanan
}

func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1)/datamakanan")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func new(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Whoa,  is near!")
}

func Read(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	makanans := make(map[int]makanan)
	rows, err := db.Query("SELECT * FROM makanan")
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {

		var id, harga int
		var nama, hidangan string

		err = rows.Scan(&id, &nama, &hidangan, &harga)
		log.Println(id)
		log.Println(nama)
		log.Println(hidangan)
		log.Println(harga)
		makanans[id] = makanan{nama, hidangan, harga}

	}
	p := view{Judul: "Data Makanan", Data: makanans}
	t, _ := template.ParseFiles("index.html")
	fmt.Println(t.Execute(w, p))
	defer db.Close()
}

func main() {
	log.Println("Server started on: http://localhost:5050")
	http.HandleFunc("/", new)
	http.HandleFunc("/data", Read)
	http.ListenAndServe(":7050", nil)
}