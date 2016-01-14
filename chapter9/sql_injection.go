package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// username := r.URL.Query().Get("username")
	username := "' OR `id`=1 OR `name`='mc"
	db, _ := sql.Open("mysql", "root@/test")

	rows, _ := db.Query("SELECT * FROM sql_injection WHERE `name`='" + username + "'")

	var id, name []byte
	for rows.Next() {
		rows.Scan(&id, &name)
		println(id, string(name))
	}
	// })
}
