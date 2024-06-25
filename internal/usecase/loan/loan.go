package loan

import (
	"context"
	"errors"
	"simple-app/internal/model"
	agrmnt "simple-app/internal/pkg/agreementLetter"
	"time"

	"log"

	"github.com/jmoiron/sqlx"
)

// Usecase instance struct for loan
type Usecase struct {
	loanRepo        loanRepo
	agreementLetter agreementLetter
}

type loanRepo interface {
	CreateTx(ctx context.Context) (*sqlx.Tx, error)

	GetByID(ctx context.Context, ID int) (loan model.Loan, err error)
	GetInvestByID(ctx context.Context, ID int) (invests []model.Invest, err error)
	GetDisburseByID(ctx context.Context, ID int) (disburses []model.Disburse, err error)
	GetList(ctx context.Context) (loan []model.Loan, err error)

	Create(ctx context.Context, param model.Loan) (data model.Loan, err error)
	Approve(ctx context.Context, param model.Approve) (id int, err error)
	Invest(ctx context.Context, dbTx *sqlx.Tx, param model.Invest) (data model.Invest, err error)
	Disburse(ctx context.Context, dbTx *sqlx.Tx, param model.Disburse) (data model.Disburse, err error)

	UpdateStatus(ctx context.Context, dbTx *sqlx.Tx, status model.Loan) (id int, err error)
}

type agreementLetter interface {
	Generate(param model.Loan) string
	Send(receiver int, agreementLetter string) error
}

// New will instantiate new loan usecase
func New(loanRepo loanRepo) *Usecase {
	return &Usecase{
		loanRepo:        loanRepo,
		agreementLetter: agrmnt.New(),
	}
}

// Create loan request, it will generate agreement letter url first, then submit it to db
func (u *Usecase) Create(ctx context.Context, param model.Loan) (data model.Loan, err error) {

	// Generate Agreement Letter
	url := u.agreementLetter.Generate(param)
	param.AgreementLetterURL = url

	data, err = u.loanRepo.Create(ctx, param)
	if err != nil {
		return data, err
	}

	return data, nil
}

// Approve loan request, it will update approver_id, status, and picture of proof url
func (u *Usecase) Approve(ctx context.Context, param model.Approve) (id int, err error) {

	loan, err := u.loanRepo.GetByID(ctx, param.ID)
	if err != nil {
		return id, err
	}

	if loan.Status >= model.APPROVED {
		return id, errors.New("status is already approved")
	}

	now := time.Now()
	param.ApprovalDate = &now
	param.Status = model.APPROVED
	id, err = u.loanRepo.Approve(ctx, param)
	if err != nil {
		return id, err
	}

	return id, nil
}

// Invest, the amount of money that given by investor
func (u *Usecase) Invest(ctx context.Context, param model.Invest) (id int, total float32, status string, err error) {

	// get approved loan
	loan, err := u.loanRepo.GetByID(ctx, param.LoanID)
	if err != nil {
		return id, total, status, err
	}

	if loan.Status != model.APPROVED {
		return id, total, status, errors.New("status of loan is invalid")
	}

	// get current investment
	invests, err := u.loanRepo.GetInvestByID(ctx, param.LoanID)
	if err != nil {
		return id, total, status, err
	}

	// calculate current total of invested
	var current float32 = 0
	for _, invest := range invests {
		current += invest.Amount
	}

	// init db transaction
	dbTx, _ := u.loanRepo.CreateTx(ctx)
	defer dbTx.Rollback()

	// insert investment
	data, err := u.loanRepo.Invest(ctx, dbTx, param)
	if err != nil {
		return id, total, status, err
	}

	// calculate total of investment
	total = current + param.Amount

	// check if total is equal of more than principal amount
	if total >= float32(loan.PrincipalAmount) {
		// Send  all of the investors the letter of agreement
		for _, invest := range invests {
			err = u.agreementLetter.Send(invest.InvestorID, loan.AgreementLetterURL)
			if err != nil {
				log.Printf("failed to send email, err = %v", err)
			}
		}

		// update status of loan to be "invested"
		loan.Status = model.INVESTED
		id, err = u.loanRepo.UpdateStatus(ctx, dbTx, loan)
		if err != nil {
			return id, total, status, err
		}
	}

	err = dbTx.Commit()
	if err != nil {
		return id, total, status, err
	}

	return data.ID, total, loan.Status.ToString(), err

}

// Disburse, the amount of money that will be disbursed by borrower
func (u *Usecase) Disburse(ctx context.Context, param model.Disburse) (id int, err error) {

	// get approved loan
	loan, err := u.loanRepo.GetByID(ctx, param.LoanID)
	if err != nil {
		return id, err
	}

	if loan.Status != model.INVESTED {
		return id, errors.New("status of loan is invalid")
	}

	// calculate current total of invested

	// init db transaction
	dbTx, _ := u.loanRepo.CreateTx(ctx)
	defer dbTx.Rollback()

	// insert investment
	now := time.Now()
	param.DisbursementDate = &now
	data, err := u.loanRepo.Disburse(ctx, dbTx, param)
	if err != nil {
		return id, err
	}

	// update status of loan to be "disbursed"
	loan.Status = model.DISBURSED
	id, err = u.loanRepo.UpdateStatus(ctx, dbTx, loan)
	if err != nil {
		return id, err
	}

	err = dbTx.Commit()
	if err != nil {
		return id, err
	}

	return data.ID, err

}

func (u *Usecase) GetDetail(ctx context.Context, id int) (detail model.Detail, err error) {

	// get loan detail
	loan, err := u.loanRepo.GetByID(ctx, id)
	if err != nil {
		return detail, err
	}

	// get investor list
	investors, err := u.loanRepo.GetInvestByID(ctx, id)
	if err != nil {
		return detail, err
	}

	// get disbursement list
	disbursement, err := u.loanRepo.GetDisburseByID(ctx, id)
	if err != nil {
		return detail, err
	}

	loan.StatusStr = loan.Status.ToString()

	// collect into one struct
	detail.Loan = loan
	detail.Investors = investors
	detail.Disbursements = disbursement

	return detail, err
}

// GetList get list of loans
func (u *Usecase) GetList(ctx context.Context) (loans []model.Loan, err error) {
	loans, err = u.loanRepo.GetList(ctx)
	if err != nil {
		return loans, err
	}

	for i, loan := range loans {
		loans[i].StatusStr = loan.Status.ToString()
	}

	return loans, nil
}
