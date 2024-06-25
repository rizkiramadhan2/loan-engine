package handler

import (
	"simple-app/app"
	"simple-app/internal/pkg/validate"
)

// Handler is struct for blog http handler
type Handler struct {
	validator *validate.Validate
	loan      app.LoanUseCase
}

// New will instantiate http blog package
func New(loanUc app.LoanUseCase) *Handler {
	v := validate.New(
		&validate.Options{})

	return &Handler{
		validator: v,
		loan:      loanUc,
	}
}
