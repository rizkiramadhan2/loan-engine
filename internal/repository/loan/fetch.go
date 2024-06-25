package loan

import (
	"database/sql"
	"errors"
	"simple-app/internal/model"

	"golang.org/x/net/context"
)

// GetByID get loan by ID
func (u *Loan) GetByID(ctx context.Context, ID int) (loan model.Loan, err error) {
	var getQuery = `SELECT id ,borrower_id ,principal_amount ,rate ,roi ,status ,agreement_letter_url ,picture_proof_url ,approver_id ,approval_date FROM loan WHERE id=$1`

	err = u.db.GetContext(ctx, &loan, getQuery, ID)
	if err != nil && err != sql.ErrNoRows {
		return loan, err
	}

	if err == sql.ErrNoRows {
		return loan, errors.New("not found")
	}

	return loan, nil
}

// GetList get list of loat
func (u *Loan) GetList(ctx context.Context) (loan []model.Loan, err error) {
	var getQuery = `SELECT id ,borrower_id ,principal_amount ,rate ,roi ,status ,agreement_letter_url ,picture_proof_url ,approver_id ,approval_date FROM loan`

	err = u.db.SelectContext(ctx, &loan, getQuery)
	if err != nil && err != sql.ErrNoRows {
		return loan, err
	}

	if err == sql.ErrNoRows {
		return loan, errors.New("not found")
	}

	return loan, nil
}

// GetInvestByID get investment by ID
func (u *Loan) GetInvestByID(ctx context.Context, ID int) (invests []model.Invest, err error) {
	var getQuery = `SELECT id ,loan_id ,investor_id ,amount FROM loan_investment WHERE loan_id=$1`

	err = u.db.SelectContext(ctx, &invests, getQuery, ID)
	if err != nil && err != sql.ErrNoRows {
		return invests, err
	}

	if err == sql.ErrNoRows {
		return invests, errors.New("not found")
	}

	return invests, nil
}

// GetDisburseByID get disbursement by ID
func (u *Loan) GetDisburseByID(ctx context.Context, ID int) (disburses []model.Disburse, err error) {
	var getQuery = `SELECT id ,loan_id ,signed_agreement_url ,disburser_employee_id, Disbursement_date FROM loan_disbursement WHERE loan_id=$1`

	err = u.db.SelectContext(ctx, &disburses, getQuery, ID)
	if err != nil && err != sql.ErrNoRows {
		return disburses, err
	}

	if err == sql.ErrNoRows {
		return disburses, errors.New("not found")
	}

	return disburses, nil
}
