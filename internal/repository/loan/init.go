package loan

import (
	"simple-app/internal/pkg/sqldb"
)

type Loan struct {
	db *sqldb.DB
}

type Param struct {
	DB *sqldb.DB
}

func New(p Param) Loan {
	return Loan{
		db: p.DB,
	}
}
