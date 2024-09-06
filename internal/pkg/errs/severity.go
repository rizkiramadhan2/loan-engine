package errs

type Severity int

const (
	SeverityUnknown Severity = iota
	SeverityWarning
	SeverityError
	SeverityFatal
)
