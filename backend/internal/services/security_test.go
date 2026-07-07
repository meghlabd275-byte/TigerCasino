package services

import (
	"fmt"
	"testing"
)

func TestSecurityBridge(t *testing.T) {
	bridge := NewSecurityBridge()
	password := "tiger123"
	hash := bridge.HashPassword(password)
	if hash == "" {
		t.Fatal("Hash is empty")
	}
	if !bridge.VerifyPassword(password, hash) {
		t.Fatal("Password verification failed")
	}
	fmt.Println("Security Bridge Passed")
}
