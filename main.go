package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/lib/pq"
)

func main() {
	hostname := flag.String("hostname", "", "hostname")
	username := flag.String("username", "", "username")
	password := flag.String("password", "", "password")
	schema := flag.String("schema", "", "schema")
	port := flag.String("port", "", "port")
	flag.Parse()

	//db, err := sql.Open("postgres", fmt.Sprintf("host=%v user=%v dbname=%v password=%v port=%v sslmode=disable", hostname, username, schema, password, port))
	conn, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", *username, *password, *hostname, *port, *schema))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	err = conn.Ping()
	handleErr(err)

	type Column struct {
		ColumnName string
		ColumnType string
	}
	type Table struct {
		TableName string
		Cols      []Column
	}

	var data [3]Table
	var tableID = 0
	tables, err := conn.Query(fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%v' ORDER BY table_name DESC;", *schema))
	handleErr(err)
	for tables.Next() {
		var tableName string
		err = tables.Scan(&tableName)
		handleErr(err)
		data[tableID].TableName = tableName
		col := make([]Column, 0)
		columns, err := conn.Query(fmt.Sprintf("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '%v' AND table_schema = '%v';", tableName, *schema))

		for columns.Next() {
			var colName string
			var colType string
			err = columns.Scan(&colName, &colType)
			handleErr(err)
			col = append(col, Column{colName, colType})
		}
		data[tableID].Cols = col
		tableID++
	}
	t, err := template.ParseFiles("struct.tmpl")
	handleErr(err)
	t.Execute(os.Stdout, data)

}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
