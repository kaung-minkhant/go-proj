package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	account, err := NewAccount("polar", "kaung", "shunnn")
	assert.Nil(t, err)
	fmt.Println("Testing creating new account", account)
}
