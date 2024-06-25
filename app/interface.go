package app

import (
	"context"

	"simple-app/internal/model"
)

type LoanUseCase interface {
	GetDetail(ctx context.Context, id int) (detail model.Detail, err error)
	GetList(ctx context.Context) (loans []model.Loan, err error)

	Create(ctx context.Context, param model.Loan) (data model.Loan, err error)
	Approve(ctx context.Context, param model.Approve) (id int, err error)
	Invest(ctx context.Context, param model.Invest) (id int, total float32, status string, err error)
	Disburse(ctx context.Context, param model.Disburse) (id int, err error)
}
