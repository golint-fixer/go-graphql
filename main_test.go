package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}
func TestConvertType(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "string", convertType("time"), "should be string")
	assert.Equal(t, "string", convertType("datetime"), "should be string")
	assert.Equal(t, "string", convertType("char"), "should be string")
	assert.Equal(t, "string", convertType("varchar"), "should be string")
	assert.Equal(t, "string", convertType("blob"), "should be string")
	assert.Equal(t, "int", convertType("integer"), "should be integer")
	assert.Equal(t, "int", convertType("int"), "should be integer")
	assert.Equal(t, "int", convertType("timestamp"), "should be integer")
	assert.Equal(t, "float", convertType("float"), "should be float")
	assert.Equal(t, "float", convertType("double"), "should be float")
	assert.Equal(t, "bool", convertType("boolean"), "should be bool")
	assert.Equal(t, "string", convertType("enum"), "should be string")
	assert.Equal(t, "string", convertType("other"), "should be string")
}
func TestTormatColName(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "Id", formatColName("id"), "should be Id")
	assert.Equal(t, "UserName", formatColName("user_name"), "should be UserName")
	assert.Equal(t, "Username", formatColName("username"), "should be Username")
	assert.Equal(t, "Col5", formatColName("col_5"), "should be Col5")
	assert.Equal(t, "Code", formatColName("code"), "should be Code")
}
func TestGetTableInfo(t *testing.T) {
	t.Parallel()
	var expected []table
	colOne := column{"Id", "int", "id"}
	colTwo := column{"Name", "string", "name"}
	expected = append(expected, table{"Tableone", []column{colOne, colTwo}})
	expected = append(expected, table{"Tabletwo", []column{colOne, colTwo}})

	// Open new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	// columns to be used for result
	tableRows := sqlmock.NewRows([]string{"table_name"}).
		AddRow("tableone").
		AddRow("tabletwo")
	colRows := sqlmock.NewRows([]string{"COLUMN_NAME", "DATA_TYPE"}).
		AddRow("id", "int").
		AddRow("name", "varchar")
	// you cant reuse mocked rows
	colRowsTwo := sqlmock.NewRows([]string{"COLUMN_NAME", "DATA_TYPE"}).
		AddRow("id", "int").
		AddRow("name", "varchar")
	mock.ExpectQuery("SELECT table_name FROM information_schema.tables").WillReturnRows(tableRows)
	mock.ExpectQuery("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE").WillReturnRows(colRows)
	mock.ExpectQuery("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE").WillReturnRows(colRowsTwo)
	data := getTableInfo(db, "some_schema")
	assert.EqualValues(t, expected, data, "should be equal")
}
func TestHandleError(t *testing.T) {
	t.Parallel()
	assert.Panics(t, func() {
		handleErr(errors.New("some error"))
	}, "Calling handleErr() should panic")
}
