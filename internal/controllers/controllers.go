package controllers

import (
	"github/billing-engine-main/internal/models"
	"github/billing-engine-main/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *services.BillingService
}

func NewController(service *services.BillingService) *Controller {
	return &Controller{service: service}
}

func (h *Controller) CreateLoan(c *gin.Context) {
	var req models.CreateLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loan, err := h.service.CreateLoan(req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, loan)
}

func (h *Controller) MakePayment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan ID"})
		return
	}

	var req models.MakePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.MakePayment(uint(id), req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment successful"})
}

func (h *Controller) GetOutstanding(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan ID"})
		return
	}

	outstanding, err := h.service.GetOutstanding(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	loanDetails, err := h.service.GetLoan(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "loan not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"1_outstanding": outstanding,
		"2_loanDetails": loanDetails,
	})
}

func (h *Controller) GetDelinquencyStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan ID"})
		return
	}

	status, err := h.service.GetLoanStatus(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}
