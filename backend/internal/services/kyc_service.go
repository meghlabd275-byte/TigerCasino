package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// KYCService handles identity verification
type KYCService struct {
	mu           sync.RWMutex
	submissions  map[string]*KYCSubmission
	users        map[string]*KYCStatus
}

// KYCSubmission represents a KYC application
type KYCSubmission struct {
	ID            string
	UserID        string
	Status        string // pending, review, approved, rejected
	Level         int    // 1: Basic, 2: Intermediate, 3: Full
	Documents     []Document
	Selfie        string
	AddressProof  string
	SubmittedAt   time.Time
	ReviewedAt    *time.Time
	RejectReason  string
}

// Document represents an uploaded document
type Document struct {
	Type     string // id_card, passport, drivers_license
	Number   string
	FrontURL string
	BackURL  string
	Expiry   string
}

// KYCStatus represents user's KYC status
type KYCStatus struct {
	UserID         string
	Level          int
	Status         string
	VerifiedAt     *time.Time
	ExpiresAt      *time.Time
	Limits         UserLimits
	Requirements   []string
}

// UserLimits represents withdrawal and deposit limits
type UserLimits struct {
	DailyDepositLimit    float64
	MonthlyDepositLimit float64
	DailyWithdrawalLimit float64
	MonthlyWithdrawalLimit float64
}

// NewKYCService creates a new KYC service
func NewKYCService() *KYCService {
	return &KYCService{
		submissions: make(map[string]*KYCSubmission),
		users:       make(map[string]*KYCStatus),
	}
}

// SubmitKYC creates a new KYC submission
func (s *KYCService) SubmitKYC(userID string, level int, documents []Document, selfie, addressProof string) (*KYCSubmission, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	submission := &KYCSubmission{
		ID:           uuid.New().String(),
		UserID:       userID,
		Status:       "pending",
		Level:        level,
		Documents:    documents,
		Selfie:       selfie,
		AddressProof: addressProof,
		SubmittedAt:  time.Now(),
	}

	s.submissions[submission.ID] = submission
	return submission, nil
}

// GetSubmission returns a KYC submission
func (s *KYCService) GetSubmission(submissionID string) (*KYCSubmission, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sub, ok := s.submissions[submissionID]
	if !ok {
		return nil, fmt.Errorf("submission not found")
	}
	return sub, nil
}

// GetUserStatus returns KYC status for a user
func (s *KYCService) GetUserStatus(userID string) (*KYCStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status, ok := s.users[userID]
	if !ok {
		// Return default unverified status
		return &KYCStatus{
			UserID:       userID,
			Level:        0,
			Status:        "unverified",
			Limits:       UserLimits{DailyDepositLimit: 1000, MonthlyDepositLimit: 5000},
			Requirements: []string{"Submit ID document", "Take selfie", "Provide address proof"},
		}, nil
	}
	return status, nil
}

// ApproveKYC approves a KYC submission
func (s *KYCService) ApproveKYC(submissionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sub, ok := s.submissions[submissionID]
	if !ok {
		return fmt.Errorf("submission not found")
	}

	sub.Status = "approved"
	now := time.Now()
	sub.ReviewedAt = &now

	// Update user status
	limits := getLimitsForLevel(sub.Level)
	s.users[sub.UserID] = &KYCStatus{
		UserID:       sub.UserID,
		Level:        sub.Level,
		Status:       "verified",
		VerifiedAt:   &now,
		ExpiresAt:    getExpiryDate(now),
		Limits:       limits,
		Requirements: []string{},
	}

	return nil
}

// RejectKYC rejects a KYC submission
func (s *KYCService) RejectKYC(submissionID, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sub, ok := s.submissions[submissionID]
	if !ok {
		return fmt.Errorf("submission not found")
	}

	sub.Status = "rejected"
	sub.RejectReason = reason
	now := time.Now()
	sub.ReviewedAt = &now

	// Update user status
	s.users[sub.UserID] = &KYCStatus{
		UserID:       sub.UserID,
		Level:        0,
		Status:       "rejected",
		Requirements: []string{reason},
	}

	return nil
}

// GetRequiredDocuments returns required documents for a level
func (s *KYCService) GetRequiredDocuments(level int) []string {
	switch level {
	case 1:
		return []string{"ID Document (Passport, National ID, or Driver's License)"}
	case 2:
		return []string{"ID Document", "Selfie with ID"}
	case 3:
		return []string{"ID Document", "Selfie with ID", "Proof of Address (Bank Statement or Utility Bill)"}
	default:
		return []string{}
	}
}

func getLimitsForLevel(level int) UserLimits {
	switch level {
	case 1:
		return UserLimits{
			DailyDepositLimit:       1000,
			MonthlyDepositLimit:     5000,
			DailyWithdrawalLimit:   500,
			MonthlyWithdrawalLimit: 2000,
		}
	case 2:
		return UserLimits{
			DailyDepositLimit:       10000,
			MonthlyDepositLimit:    50000,
			DailyWithdrawalLimit:   5000,
			MonthlyWithdrawalLimit: 25000,
		}
	case 3:
		return UserLimits{
			DailyDepositLimit:       100000,
			MonthlyDepositLimit:    500000,
			DailyWithdrawalLimit:   50000,
			MonthlyWithdrawalLimit: 250000,
		}
	default:
		return UserLimits{
			DailyDepositLimit:       1000,
			MonthlyDepositLimit:     5000,
			DailyWithdrawalLimit:   500,
			MonthlyWithdrawalLimit: 2000,
		}
	}
}

func getExpiryDate(verifiedAt time.Time) *time.Time {
	expiry := verifiedAt.AddDate(1, 0, 0) // 1 year validity
	return &expiry
}

// CalculateRemainingLimits calculates remaining limits for a user
func (s *KYCService) CalculateRemainingLimits(userID, period string, deposited, withdrawn float64) (float64, float64, error) {
	status, err := s.GetUserStatus(userID)
	if err != nil {
		return 0, 0, err
	}

	var dailyDeposit, dailyWithdrawal float64
	if period == "daily" {
		dailyDeposit = status.Limits.DailyDepositLimit - deposited
		dailyWithdrawal = status.Limits.DailyWithdrawalLimit - withdrawn
	} else {
		dailyDeposit = status.Limits.MonthlyDepositLimit - deposited
		dailyWithdrawal = status.Limits.MonthlyWithdrawalLimit - withdrawn
	}

	if dailyDeposit < 0 {
		dailyDeposit = 0
	}
	if dailyWithdrawal < 0 {
		dailyWithdrawal = 0
	}

	return dailyDeposit, dailyWithdrawal, nil
}
