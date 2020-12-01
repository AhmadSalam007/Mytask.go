package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	csvreader()
}
func csvreader() {
	// Open the file
	recordfile, err := os.Open("./Mytask.csv")
	if err != nil {
		fmt.Println("Couldn't open the csv file", err)
	}
	r := csv.NewReader(recordfile)
	Mytask, _ := r.ReadAll()
	fmt.Println(Mytask)
}

type person struct {
	firstname  string
	lastname   string
	age        int
	bloodgroup string
}

func dbconn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "Ali202271"
	dbName := "task"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

//Index is used to access the controller
func Index(w http.ResponseWriter, r *http.Request) {
	db := dbconn()
	selDB, err := db.Query("SELECT * FROM person ORDER BY firstname DESC")
	if err != nil {
		panic(err.Error())
	}
	per := person{}
	res := []person{}
	for selDB.Next() {
		var firstname, lastname, bloodgroup string
		var age int
		err = selDB.Scan(&firstname, &lastname, &age, &bloodgroup)
		if err != nil {
			panic(err.Error())
		}
		per.firstname = firstname
		per.lastname = lastname
		per.age = age
		per.bloodgroup = bloodgroup
		res = append(res, per)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

// Show is used to show the selected variable
func Show(w http.ResponseWriter, r *http.Request) {
	db := dbconn()
	nper := r.URL.Query().Get("firstname")
	selDB, err := db.Query("SELECT * FROM person WHERE firstname=?", nper)
	if err != nil {
		panic(err.Error())
	}
	per := person{}
	for selDB.Next() {
		var age int
		var firstname, lastname, bloodgroup string
		err = selDB.Scan(&firstname, &lastname, &age, &bloodgroup)
		if err != nil {
			panic(err.Error())
		}
		per.firstname = firstname
		per.lastname = lastname
		per.age = age
		per.bloodgroup = bloodgroup
	}
	tmpl.ExecuteTemplate(w, "Show", per)
	defer db.Close()
}

//New is used to enter new variable
func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

//Insert is used to submit the form
func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbconn()
	if r.Method == "POST" {
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lasttname")
		age := r.FormValue("age")
		bloodgroup := r.FormValue("bloodgroup")
		insForm, err := db.Prepare("INSERT INTO person(firstname, lastname, age, bloodgroup) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(firstname, lastname, age, bloodgroup)
		log.Println("INSERT: fisrtname: " + firstname + " | lastname: " + lastname + " | age: " + age + " | bloodgroup: " + bloodgroup)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}
func handle() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/insert", Insert)
	http.ListenAndServe(":8080", nil)
}
