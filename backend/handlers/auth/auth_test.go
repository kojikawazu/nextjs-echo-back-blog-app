package handlers_auth

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	SetupTest(nil)

	code := m.Run()
	os.Exit(code)
}
