package agreementletter

import (
	"simple-app/internal/model"

	"log"
)

type AgreementLetter struct{}

func New() *AgreementLetter {
	return &AgreementLetter{}
}

func (a *AgreementLetter) Generate(model.Loan) string {
	/* THIS IS JUST PLACEHOLDER TO GENERATE URL OF AGREEMENT LETTER */
	return "http://example-of-agreement-letter.com"
}

func (a *AgreementLetter) Send(receiver int, agreementLetter string) error {
	/* THIS IS JUST PLACEHOLDER TO SEND EMAIL OF AGREEMENT LETTER */
	log.Printf("sending agreement letter to %v", receiver)
	return nil
}
