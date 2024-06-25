package model

import "time"

type Loan struct {
	ID                 int        `json:"id" db:"id"`
	BorrowerID         int        `json:"borrower_id" db:"borrower_id" validate:"required"`
	PrincipalAmount    float64    `json:"principal_amount" db:"principal_amount" validate:"required"`
	Rate               float64    `json:"rate" db:"rate" validate:"required"`
	Roi                float64    `json:"roi" db:"roi" validate:"required"`
	Status             LoanStatus `json:"status" db:"status"`
	StatusStr          string     `json:"status_str,omitempty"`
	AgreementLetterURL string     `json:"agreement_letter_url" db:"agreement_letter_url"`
	PictureProofURL    *string    `json:"picture_proof_url,omitempty" db:"picture_proof_url"`
	ApproverID         *int       `json:"approver_id,omitempty" db:"approver_id"`
	ApprovalDate       *time.Time `json:"approval_date,omitempty" db:"approval_date"`
}

type Approve struct {
	ID              int        `json:"id" db:"id" validate:"required"`
	PictureProofURL *string    `json:"picture_proof_url" db:"picture_proof_url" validate:"required"`
	ApproverID      int        `json:"approver_id" db:"approver_id"  validate:"required"`
	ApprovalDate    *time.Time `json:"approval_date" db:"approval_date"`
	Status          LoanStatus `json:"status" db:"status"`
}

type Invest struct {
	ID         int        `json:"id"`
	LoanID     int        `json:"loan_id" db:"loan_id" validate:"required"`
	InvestorID int        `json:"investor_id" db:"investor_id"  validate:"required"`
	Amount     float32    `json:"amount" db:"amount"  validate:"required"`
	Status     LoanStatus `json:"status,omitempty"`
}

type Disburse struct {
	ID                 int        `json:"id"`
	LoanID             int        `json:"loan_id" db:"loan_id" validate:"required"`
	SignedAgreementURL string     `json:"signed_agreement_url"   db:"signed_agreement_url" validate:"required"`
	DisburseEmployeeID int        `json:"disburser_employee_id"  db:"disburser_employee_id" validate:"required"`
	DisbursementDate   *time.Time `json:"disbursement_date"  db:"disbursement_date"`
}

type Detail struct {
	Loan
	Investors     []Invest   `json:"investors"`
	Disbursements []Disburse `json:"disbursement"`
}

type LoanStatus int

const (
	PROPOSED LoanStatus = iota + 1
	APPROVED
	INVESTED
	DISBURSED
)

func (ls LoanStatus) ToString() string {
	switch true {
	case ls == PROPOSED:
		return "proposed"
	case ls == APPROVED:
		return "approved"
	case ls == INVESTED:
		return "invested"
	case ls == DISBURSED:
		return "disbursed"
	}

	return "invalid"
}
