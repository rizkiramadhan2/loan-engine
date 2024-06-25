package loan

import (
	"database/sql"
	"fmt"
	"simple-app/internal/model"

	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
)

// CreateTx start transaction if needed
func (l *Loan) CreateTx(ctx context.Context) (*sqlx.Tx, error) {
	return l.db.GetMaster().BeginTxx(ctx, &sql.TxOptions{})
}

// Create create loan
func (l *Loan) Create(ctx context.Context, param model.Loan) (data model.Loan, err error) {
	// Insert file information into the database
	// STATUS DEFAULT IS "proposed"
	query := `
		INSERT INTO loan (
			borrower_id,
			principal_amount,
			rate,
			roi,
			agreement_letter_url
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		) RETURNING
		 	id,
			status,
			borrower_id,
			principal_amount,
			rate,
			roi,
			agreement_letter_url
	`

	err = l.db.GetContext(ctx, &data, query,
		param.BorrowerID,
		param.PrincipalAmount,
		param.Rate,
		param.Roi,
		param.AgreementLetterURL,
	)
	if err != nil {
		return data, fmt.Errorf("failed to insert loan: %w", err)
	}

	return data, nil
}

func (l *Loan) Approve(ctx context.Context, param model.Approve) (id int, err error) {
	q := `
		UPDATE loan SET
			picture_proof_url = :picture_proof_url,
			approver_id = :approver_id,
			approval_date = :approval_date,
			status = :status
		WHERE id = :id
		RETURNING
			id
	`

	q, arg, err := sqlx.Named(q, param)
	if err != nil {
		return id, err
	}

	err = l.db.GetContext(ctx, &id, l.db.Rebind(q), arg...)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (l *Loan) UpdateStatus(ctx context.Context, dbTx *sqlx.Tx, status model.Loan) (id int, err error) {
	querier := dbTx
	if dbTx == nil {
		querier = l.db.GetMaster().MustBegin()
	}

	q := `
	UPDATE loan SET
		status = :status
	WHERE id = :id
	RETURNING
		id
`

	q, arg, err := sqlx.Named(q, status)
	if err != nil {
		return id, err
	}

	err = querier.GetContext(ctx, &id, l.db.Rebind(q), arg...)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (l *Loan) Invest(ctx context.Context, dbTx *sqlx.Tx, param model.Invest) (data model.Invest, err error) {

	querier := dbTx
	if dbTx == nil {
		querier = l.db.GetMaster().MustBegin()
	}

	query := `
		INSERT INTO loan_investment (
			loan_id,
			investor_id,
			amount
		) VALUES (
			$1,
			$2,
			$3
		) RETURNING
		 	id,
			loan_id
			investor_id,
			amount
	`

	err = querier.GetContext(ctx, &data, query,
		param.LoanID,
		param.InvestorID,
		param.Amount,
	)
	if err != nil {
		return data, fmt.Errorf("failed to invest loan: %w", err)
	}

	return data, nil
}

func (l *Loan) Disburse(ctx context.Context, dbTx *sqlx.Tx, param model.Disburse) (data model.Disburse, err error) {

	querier := dbTx
	if dbTx == nil {
		querier = l.db.GetMaster().MustBegin()
	}

	query := `
		INSERT INTO loan_disbursement (
			loan_id,
			signed_agreement_url,
			disburser_employee_id,
			disbursement_date
			
		) VALUES (
			$1,
			$2,
			$3,
			$4
		) RETURNING
		 	id,
			loan_id,
			signed_agreement_url,
			disburser_employee_id,
			disbursement_date
	`

	err = querier.GetContext(ctx, &data, query,
		param.LoanID,
		param.SignedAgreementURL,
		param.DisburseEmployeeID,
		param.DisbursementDate,
	)
	if err != nil {
		return data, fmt.Errorf("failed to disburse loan: %w", err)
	}

	return data, nil
}
