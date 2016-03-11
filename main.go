package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"unicode"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/lib/pq"
)

func main() {
	hostname := flag.String("hostname", "", "hostname")
	username := flag.String("username", "", "username")
	password := flag.String("password", "", "password")
	schema := flag.String("schema", "", "schema")
	port := flag.String("port", "", "port")
	output := flag.String("output", "", "output")
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
		ColumnName   string
		ColumnType   string
		DBColumnName string
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
		var col []Column
		columns, err := conn.Query(fmt.Sprintf("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '%v' AND table_schema = '%v';", tableName, *schema))

		for columns.Next() {
			var colName string
			var colType string
			err = columns.Scan(&colName, &colType)
			handleErr(err)
			formatColName(colName)
			col = append(col, Column{formatColName(colName), convertType(colType), colName})
		}
		data[tableID].Cols = col
		tableID++
	}

	// process template
	t, err := template.ParseFiles("struct.tmpl")
	handleErr(err)

	// create file
	f, err := os.Create(*output)
	handleErr(err)
	defer f.Close()
	t.Execute(f, data)

	// format the file
	cmd := exec.Command("gofmt", "-w", f.Name())
	err = cmd.Run()
	handleErr(err)
}

// formatColName formats the column name into camel case
func formatColName(name string) string {
	arr := []byte(name)
	var output []byte
	capNextChar := false
	for i, char := range arr {
		if i == 0 || capNextChar {
			output = append(output, byte(unicode.ToUpper(rune(char))))
			capNextChar = false
		} else if char == '_' {
			capNextChar = true
		} else {
			output = append(output, char)
			capNextChar = false
		}
	}
	return string(output)
}

// convertType converts the db col type to the corresponding go type
func convertType(dbType string) string {
	switch dbType {
	// Dates represented as strings
	case "time", "date", "datetime":
		fallthrough

	// Buffers represented as strings
	case "bit", "blob", "tinyblob", "longblob", "mediumblob", "binary", "varbinary":
		fallthrough

	// Numbers that may exceed float precision, repesent as string
	case "bigint", "decimal", "numeric", "geometry", "bigserial":
		fallthrough

	// Network addresses represented as strings
	case "cidr", "inet", "macaddr":
		fallthrough

	// Strings
	case "set", "char", "text", "uuid", "varchar", "nvarchar", "tinytext", "longtext", "character", "mediumtext":
		return "string"
	// Integers
	case "int", "year", "serial", "integer", "tinyint", "smallint", "mediumint", "timestamp":
		return "int"
	// Floats
	case "real", "float", "double", "double precision":
		return "float"

	// Booleans
	case "boolean":
		return "bool"

	// Enum special case
	case "enum":
		return "string"

	default:
		return "string"
	}
}

// handleErr handles errors in a consistant way
func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
