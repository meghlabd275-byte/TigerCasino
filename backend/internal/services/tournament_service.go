package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TournamentService manages casino tournaments
type TournamentService struct {
	db                *gorm.DB
	activeTournaments map[string]*Tournament
	mu                sync.RWMutex
}

// Tournament represents a tournament
type Tournament struct {
	ID          string
	Name        string
	Type        string
	StartTime   time.Time
	EndTime     time.Time
	MinBet      float64
	PrizePool   float64
	Prizes      []Prize
	Status      string
	Leaderboard map[string]float64
	mu          sync.RWMutex
}

// Prize represents a tournament prize
type Prize struct {
	Rank    int
	MinRank int
	MaxRank int
	Amount  float64
}

// NewTournamentService creates a new tournament service
func NewTournamentService(db *gorm.DB) *TournamentService {
	s := &TournamentService{
		db:                db,
		activeTournaments: make(map[string]*Tournament),
	}
	s.initializeDefaultTournaments()
	return s
}

func (s *TournamentService) initializeDefaultTournaments() {
	slotsTournament := &Tournament{
		ID:          uuid.New().String(),
		Name:        "Daily Slots Tournament",
		Type:        "slots",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(24 * time.Hour),
		MinBet:      0.50,
		PrizePool:  1000.0,
		Prizes:      []Prize{{Rank: 1, MinRank: 1, MaxRank: 1, Amount: 300}, {Rank: 2, MinRank: 2, MaxRank: 2, Amount: 200}, {Rank: 3, MinRank: 3, MaxRank: 3, Amount: 100}},
		Status:      "active",
		Leaderboard: make(map[string]float64),
	}
	s.activeTournaments["daily_slots"] = slotsTournament

	blackjackTournament := &Tournament{
		ID:          uuid.New().String(),
		Name:        "Weekly Blackjack Challenge",
		Type:        "blackjack",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(7 * 24 * time.Hour),
		MinBet:      5.0,
		PrizePool:   5000.0,
		Prizes:      []Prize{{Rank: 1, MinRank: 1, MaxRank: 1, Amount: 1500}, {Rank: 2, MinRank: 2, MaxRank: 2, Amount: 1000}},
		Status:      "active",
		Leaderboard: make(map[string]float64),
	}
	s.activeTournaments["weekly_blackjack"] = blackjackTournament

	crashTournament := &Tournament{
		ID:          uuid.New().String(),
		Name:        "Monthly Crash Masters",
		Type:        "crash",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(30 * 24 * time.Hour),
		MinBet:      1.0,
		PrizePool:   10000.0,
		Prizes:      []Prize{{Rank: 1, MinRank: 1, MaxRank: 1, Amount: 3000}, {Rank: 2, MinRank: 2, MaxRank: 2, Amount: 2000}},
		Status:      "active",
		Leaderboard: make(map[string]float64),
	}
	s.activeTournaments["monthly_crash"] = crashTournament
}

// GetActiveTournaments returns all active tournaments
func (s *TournamentService) GetActiveTournaments() []Tournament {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var tournaments []Tournament
	for _, t := range s.activeTournaments {
		if t.Status == "active" {
			tournaments = append(tournaments, *t)
		}
	}
	return tournaments
}

// GetTournament returns a specific tournament
func (s *TournamentService) GetTournament(id string) (*Tournament, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tournament, ok := s.activeTournaments[id]
	if !ok {
		return nil, fmt.Errorf("tournament not found")
	}
	return tournament, nil
}

// RecordBet records a bet for the tournament leaderboard
func (s *TournamentService) RecordBet(userID, tournamentID, gameType string, betAmount float64) error {
	s.mu.RLock()
	tournament, ok := s.activeTournaments[tournamentID]
	s.mu.RUnlock()
	if !ok {
		return fmt.Errorf("tournament not found")
	}
	if tournament.Type != "overall" && tournament.Type != gameType {
		return nil
	}
	if betAmount < tournament.MinBet {
		return nil
	}
	tournament.mu.Lock()
	defer tournament.mu.Unlock()
	currentPoints := tournament.Leaderboard[userID]
	tournament.Leaderboard[userID] = currentPoints + betAmount
	return nil
}

// GetLeaderboard returns the tournament leaderboard
func (s *TournamentService) GetLeaderboard(tournamentID string) ([]map[string]interface{}, error) {
	s.mu.RLock()
	tournament, ok := s.activeTournaments[tournamentID]
	s.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("tournament not found")
	}
	tournament.mu.RLock()
	defer tournament.mu.RUnlock()
	type userPoints struct {
		userID string
		points float64
	}
	var sorted []userPoints
	for userID, points := range tournament.Leaderboard {
		sorted = append(sorted, userPoints{userID: userID, points: points})
	}
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].points > sorted[i].points {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	var leaderboard []map[string]interface{}
	for rank, up := range sorted {
		leaderboard = append(leaderboard, map[string]interface{}{
			"rank":   rank + 1,
			"user_id": up.userID,
			"points":  up.points,
		})
	}
	return leaderboard, nil
}

// GetUserRank returns a user's rank in a tournament
func (s *TournamentService) GetUserRank(tournamentID, userID string) (int, float64, error) {
	s.mu.RLock()
	tournament, ok := s.activeTournaments[tournamentID]
	s.mu.RUnlock()
	if !ok {
		return 0, 0, fmt.Errorf("tournament not found")
	}
	tournament.mu.RLock()
	defer tournament.mu.RUnlock()
	userPoints := tournament.Leaderboard[userID]
	rank := 1
	for _, points := range tournament.Leaderboard {
		if points > userPoints {
			rank++
		}
	}
	return rank, userPoints, nil
}
