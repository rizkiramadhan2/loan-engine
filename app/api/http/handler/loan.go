package handler

import (
	"errors"
	"net/http"
	"simple-app/internal/model"
	"simple-app/internal/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateLoan is a handler that validates the request
func (h *Handler) CreateLoan(c *gin.Context) {
	var loan model.Loan
	err := c.ShouldBindJSON(&loan)
	if err != nil {
		response.Err(c, response.WrapErrCode(err, response.BadRequestErrCode), err.Error())
		return
	}

	val := h.validator.ValidateStruct(loan)
	if len(val) > 0 {
		response.Err(c, response.InvalidRequestPayloadCode(val...))
		return
	}

	data, err := h.loan.Create(c.Request.Context(), loan)
	if err != nil {
		response.Err(c, response.WrapErrCode(err, response.InternalErrCode), err.Error())
		return
	}

	data.StatusStr = data.Status.ToString()

	c.JSON(http.StatusOK, data)
}

// ApproveLoan is a handler that validates the request
func (h *Handler) ApproveLoan(c *gin.Context) {
	var loan model.Approve
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.Err(c, response.WrapErrCode(errors.New("id is invalid"), response.BadRequestErrCode))
		return
	}

	loan.ID = idInt
	err = c.ShouldBindJSON(&loan)
	if err != nil {
		response.Err(c, response.WrapErrCode(err, response.BadRequestErrCode), err.Error())
		return
	}

	val := h.validator.ValidateStruct(loan)
	if len(val) > 0 {
		response.Err(c, response.InvalidRequestPayloadCode(val...))
		return
	}

	res, err := h.loan.Approve(c.Request.Context(), loan)
	if err != nil {
		if err.Error() == "status is already approved" {
			response.Err(c, response.WrapErrCode(err, response.BadRequestErrCode), err.Error())
			return
		}

		response.Err(c, response.WrapErrCode(err, response.InternalErrCode), err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"loan_id": res, "status": model.APPROVED.ToString()})
}

// InvestLoan is a handler that invest to the borrower
func (h *Handler) InvestLoan(c *gin.Context) {
	var invest model.Invest
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.Err(c, response.WrapErrCode(errors.New("id is invalid"), response.BadRequestErrCode))
		return
	}

	invest.LoanID = idInt
	err = c.ShouldBindJSON(&invest)
	if err != nil {
		response.Err(c, response.WrapErrCode(err, response.BadRequestErrCode), err.Error())
		return
	}

	val := h.validator.ValidateStruct(invest)
	if len(val) > 0 {
		response.Err(c, response.InvalidRequestPayloadCode(val...))
		return
	}

	idLoan, total, status, err := h.loan.Invest(c.Request.Context(), invest)
	if err != nil {
		if err.Error() == "status of loan is invalid" {
			response.Err(c, response.WrapErrCode(err, response.BadRequestErrCode), err.Error())
			return
		}

		response.Err(c, response.WrapErrCode(err, response.InternalErrCode), err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"loan_id": idLoan, "total_of_invested": total, "status": status})
}

// DisburseLoan is a handler that disburse by borrower
func (h *Handler) DisburseLoan(c *gin.Context) {
	var disburse model.Disburse
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.Err(c, response.WrapErrCode(errors.New("id is invalid"), response.BadRequestErrCode))
		return
	}

	disburse.LoanID = idInt
	err = c.ShouldBindJSON(&disburse)
	if err != nil {
		response.Err(c, response.WrapErrCode(err, response.BadRequestErrCode), err.Error())
		return
	}

	val := h.validator.ValidateStruct(disburse)
	if len(val) > 0 {
		response.Err(c, response.InvalidRequestPayloadCode(val...))
		return
	}

	idLoan, err := h.loan.Disburse(c.Request.Context(), disburse)
	if err != nil {
		if err.Error() == "status of loan is invalid" {
			response.Err(c, response.WrapErrCode(err, response.BadRequestErrCode), err.Error())
			return
		}

		response.Err(c, response.WrapErrCode(err, response.InternalErrCode), err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"loan_id": idLoan, "status": model.DISBURSED.ToString()})
}

// GetDetail is a handler that get the detail of loan
func (h *Handler) GetDetail(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.Err(c, response.WrapErrCode(errors.New("id is invalid"), response.BadRequestErrCode))
		return
	}

	detail, err := h.loan.GetDetail(c.Request.Context(), idInt)
	if err != nil {
		response.Err(c, response.WrapErrCode(err, response.InternalErrCode), err.Error())
		return
	}

	c.JSON(http.StatusOK, detail)
}

// GetList is a handler that get list of loan
func (h *Handler) GetList(c *gin.Context) {
	list, err := h.loan.GetList(c.Request.Context())
	if err != nil {
		response.Err(c, response.WrapErrCode(err, response.InternalErrCode), err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
}
