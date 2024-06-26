package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	type config struct {
		Name     string
		Addr     string
		Username string
		Password string
	}
	var dbConf config

	c := New(WithDir("../../test/config/"))
	if err := c.Load("database", &dbConf, nil); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, dbConf.Name, "chat")
	assert.Equal(t, dbConf.Addr, "127.0.0.1:3306")
	assert.Equal(t, dbConf.Username, "root")
	assert.Equal(t, dbConf.Password, "root")
}
