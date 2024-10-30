package main

import (
	"authentication/data"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	repo := data.NewPostgresTestRepository(nil)
	os.Exit(m.Run())
}
