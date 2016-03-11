package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {

	// assert equality
	assert.Equal(t, "string", convertType("time"), "should be string")
	assert.Equal(t, "string", convertType("datetime"), "should be string")
	assert.Equal(t, "string", convertType("char"), "should be string")
	assert.Equal(t, "string", convertType("varchar"), "should be string")
	assert.Equal(t, "string", convertType("blob"), "should be string")
	assert.Equal(t, "int", convertType("integer"), "should be integer")
	assert.Equal(t, "int", convertType("int"), "should be integer")
	assert.Equal(t, "int", convertType("timestamp"), "should be integer")
	assert.Equal(t, "bool", convertType("boolean"), "should be bool")
	assert.Equal(t, "string", convertType("enum"), "should be string")
}
