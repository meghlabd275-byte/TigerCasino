package services

/*
#cgo LDFLAGS: -L${SRCDIR}/../../../game-engine/build -ltigercasino_game_engine
double cpp_generate_dice_roll();
int cpp_calculate_slots(int bet);
int cpp_spin_roulette();
int cpp_play_blackjack(int bet);
int cpp_play_video_poker(int bet);
*/
import "C"

type GameEngineBridge struct{}

func NewGameEngineBridge() *GameEngineBridge {
	return &GameEngineBridge{}
}

func (b *GameEngineBridge) GenerateDiceRoll() float64 {
	return float64(C.cpp_generate_dice_roll())
}

func (b *GameEngineBridge) CalculateSlots(bet int) int {
	return int(C.cpp_calculate_slots(C.int(bet)))
}

func (b *GameEngineBridge) SpinRoulette() int {
	return int(C.cpp_spin_roulette())
}

func (b *GameEngineBridge) PlayBlackjack(bet int) int {
	return int(C.cpp_play_blackjack(C.int(bet)))
}

func (b *GameEngineBridge) PlayVideoPoker(bet int) int {
	return int(C.cpp_play_video_poker(C.int(bet)))
}
