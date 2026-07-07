package services

import (
	"fmt"
	"testing"
)

func TestGameEngineBridge(t *testing.T) {
	bridge := NewGameEngineBridge()
	roll := bridge.GenerateDiceRoll()
	if roll < 0 || roll > 100 {
		t.Fatalf("Dice roll out of range: %f", roll)
	}
	fmt.Println("Game Engine Bridge Passed")
}
