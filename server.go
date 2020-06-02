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

func Tambahdata(w http.ResponseWriter, r *http.Request) {
	var x string
	p := x
	t, _ := template.ParseFiles("Formtambah.html")
	fmt.Println(t.Execute(w, p))
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		nama := r.FormValue("nama")
		jenis := r.FormValue("jenis")
		harga := r.FormValue("harga")
		insForm, err := db.Prepare("INSERT INTO makanan(nama_makanan, jenis_hidangan, harga_makanan ) VALUES(?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(nama, jenis, harga)
		log.Println("INSERT: Nama: " + nama + " | Hidangan: " + jenis + " | Harga: " + harga)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM makanan WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		nama := r.FormValue("nama")
		jenis := r.FormValue("jenis")
		harga := r.FormValue("harga")
		id := r.FormValue("id")
		insForm, err := db.Prepare("UPDATE makanan SET nama_makanan=?, jenis_hidangan=?, harga_makanan=? WHERE Id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(nama, jenis, harga, id)
		log.Println("Update: Nama: " + nama + " | Hidangan: " + jenis + " | Harga: " + harga)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM makanan WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	makanans := make(map[int]makanan)
	for selDB.Next() {
		var id, harga int
		var nama, hidangan string
		err = selDB.Scan(&id, &nama, &hidangan, &harga)
		if err != nil {
			panic(err.Error())
		}
		makanans[id] = makanan{nama, hidangan, harga}
	}
	p := view{Judul: "Edit Data", Data: makanans}
	t, _ := template.ParseFiles("Formedit.html")
	fmt.Println(t.Execute(w, p))
	defer db.Close()
}

func main() {
	log.Println("Server started on: http://localhost:7050")
	http.HandleFunc("/", Read)
	http.HandleFunc("/add", Insert)
	http.HandleFunc("/tambahdata", Tambahdata)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)
	http.ListenAndServe(":7050", nil)
}
