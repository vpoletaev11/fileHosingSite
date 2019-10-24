package tests

import (
	"testing"

	"github.com/vpoletaev11/fileHostingSite/pages/login"
)

// TestLoginPage tests login page (your cap)
func TestLoginPage(t *testing.T) {
	url := "http://127.0.0.1"
	handler := login.Page(db)
}
