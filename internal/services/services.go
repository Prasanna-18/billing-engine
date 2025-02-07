package services

import (
	"errors"
	"github/billing-engine-main/internal/models"
	"time"

	"gorm.io/gorm"
)

type BillingService struct {
	db *gorm.DB
}

func NewBillingService(db *gorm.DB) *BillingService {
	return &BillingService{db: db}
}

func (s *BillingService) CreateLoan(amount float64) (*models.Loan, error) {
	if amount != 5000000 {
		return nil, errors.New("only loans of Rp 5,000,000 are supported")
	}

	loan := &models.Loan{
		Amount:    amount,
		Interest:  10.0,
		StartDate: time.Now(),
		Status:    "ACTIVE",
	}

	tx := s.db.Begin()

	if err := tx.Create(loan).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Calculate weekly payment
	weeklyPayment := (amount * (1 + 0.1)) / 50

	// Create schedule for 50 weeks
	for week := 1; week <= 50; week++ {
		schedule := &models.Schedule{
			LoanID:     loan.ID,
			WeekNumber: week,
			DueAmount:  weeklyPayment,
			DueDate:    loan.StartDate.AddDate(0, 0, week*7),
			Status:     "PENDING",
		}
		if err := tx.Create(schedule).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return loan, nil
}

func (s *BillingService) GetLoan(id uint) (*models.Loan, error) {
	var loan models.Loan
	if err := s.db.Preload("Schedules").First(&loan, id).Error; err != nil {
		return nil, err
	}
	return &loan, nil
}

func (s *BillingService) GetOutstanding(loanID uint) (float64, error) {
	var outstanding float64
	err := s.db.Model(&models.Schedule{}).
		Where("loan_id = ? AND status = ?", loanID, "PENDING").
		Select("SUM(due_amount)").
		Scan(&outstanding).Error
	return outstanding, err
}

func (s *BillingService) IsDelinquent(loanID uint) (bool, error) {
	var consecutiveMissed int
	currentTime := time.Now()

	var schedules []models.Schedule
	err := s.db.Where("loan_id = ? AND status = ? AND due_date < ?",
		loanID, "PENDING", currentTime).
		Order("week_number asc").
		Find(&schedules).Error

	if err != nil {
		return false, err
	}

	for _, schedule := range schedules {
		if schedule.Status == "PENDING" {
			consecutiveMissed++
			if consecutiveMissed >= 2 {
				return true, nil
			}
		} else {
			consecutiveMissed = 0
		}
	}

	return false, nil
}

func (s *BillingService) MakePayment(loanID uint, amount float64) error {
	var schedule models.Schedule
	err := s.db.Where("loan_id = ? AND status = ?", loanID, "PENDING").
		Order("week_number asc").
		First(&schedule).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("no pending payments found")
		}
		return err
	}

	if amount != schedule.DueAmount {
		return errors.New("payment amount must match the scheduled amount")
	}

	sevenDaysLater := time.Now().Add(7 * 24 * time.Hour)
	if schedule.DueDate.After(sevenDaysLater) {
		return errors.New("payment for the current week has already been made") // no pending due for current week
	}

	tx := s.db.Begin()

	payment := &models.Payment{
		LoanID:      loanID,
		ScheduleID:  schedule.ID,
		Amount:      amount,
		PaymentDate: time.Now(),
	}

	if err := tx.Create(payment).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&schedule).Update("status", "PAID").Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *BillingService) GetLoanStatus(loanID uint) (*models.LoanStatus, error) {
	outstanding, err := s.GetOutstanding(loanID)
	if err != nil {
		return nil, err
	}

	isDelinquent, err := s.IsDelinquent(loanID)
	if err != nil {
		return nil, err
	}

	return &models.LoanStatus{
		Outstanding:  outstanding,
		IsDelinquent: isDelinquent,
	}, nil
}
