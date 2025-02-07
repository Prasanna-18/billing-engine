package models

import (
	"time"
	"gorm.io/gorm"
)

type Loan struct {
	gorm.Model
	Amount    float64    `json:"amount"`
	Interest  float64    `json:"interest"`
	StartDate time.Time  `json:"start_date"`
	Status    string     `json:"status"`
	Schedules []Schedule `json:"schedules,omitempty"`
}

type Schedule struct {
	gorm.Model
	LoanID     uint      `json:"loan_id"`
	WeekNumber int       `json:"week_number"`
	DueAmount  float64   `json:"due_amount"`
	DueDate    time.Time `json:"due_date"`
	Status     string    `json:"status"`
	Loan       Loan      `json:"-" gorm:"foreignKey:LoanID"`
	Payments   []Payment `json:"payments,omitempty"`
}

type Payment struct {
	gorm.Model
	LoanID      uint      `json:"loan_id"`
	ScheduleID  uint      `json:"schedule_id"`
	Amount      float64   `json:"amount"`
	PaymentDate time.Time `json:"payment_date"`
	Loan        Loan      `json:"-" gorm:"foreignKey:LoanID"`
	Schedule    Schedule  `json:"-" gorm:"foreignKey:ScheduleID"`
}

type LoanStatus struct {
	Outstanding  float64 `json:"outstanding"`
	IsDelinquent bool    `json:"is_delinquent"`
}

type CreateLoanRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

type MakePaymentRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}