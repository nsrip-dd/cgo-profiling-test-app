package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/nsrip-dd/cgotraceback"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

const schema = `
create table if not exists colors (
	name text primary key,
	color text
);
PRAGMA journal_mode=WAL;
`

const update = `
insert into colors(name, color) values (?, ?)
on conflict(name) do update set color = excluded.color;
`

const query = `select color from colors where name = ?;`

func main() {
	if err := profiler.Start(); err != nil {
		log.Fatal("starting profile:", err)
	}
	defer profiler.Stop()
	addr := "localhost:8765"
	if v := os.Getenv("ADDR"); v != "" {
		addr = v
	}

	db, err := sql.Open("sqlite3", "colors.db")
	if err != nil {
		log.Fatal("opening:", err)
	}
	defer db.Close()

	if _, err := db.Exec(schema); err != nil {
		log.Fatal("creating:", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		v := r.URL.Query()
		name, color := v.Get("name"), v.Get("color")
		if name == "" || color == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tx, err := db.Begin()
		if err != nil {
			log.Printf("begin tx failed: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		stmt, err := tx.Prepare(update)
		if err != nil {
			log.Printf("prepare tx failed: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer stmt.Close()
		if _, err = stmt.Exec(name, color); err != nil {
			log.Printf("stmt exec failed: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := tx.Commit(); err != nil {
			log.Printf("stmt commit failed: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		result, err := db.Query(query, name)
		if err != nil {
			log.Printf("query failed: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var color string
		for result.Next() {
			result.Scan(&color)
		}
		if color == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "%s's favorite color is %s\n", name, color)
	})
	http.ListenAndServe(addr, mux)
}
