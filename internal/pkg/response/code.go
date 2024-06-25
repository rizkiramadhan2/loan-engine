package response

import "net/http"

const (
	genericErrorMsg        = "Uh oh! Something unexpected occurred. Please send a report to <support contact> so we can fix it."
	paymentRequiredMsg     = "You need to upgrade your plan to access this feature."
	forbiddenMsg           = "You do not have the necessary permissions to view this item."
	notFoundMsg            = "Oops! We couldn't find what you were looking for. Please send a report to <support contact> if you believe this was an error."
	unprocessableEntityMsg = "The request could not be processed correctly due to a mistake in the information provided. Please review and try again."
	tooManyRequestMsg      = "You've made too many requests. Please take a break and try again later."
	unauthorizedMsg        = "Sorry, you need to be logged in to access this page. Please log in and try again."
)

// Code type
type Code struct {
	code     string
	httpCode int
	devMsg   string
	userMsg  string
}

// NewCode constructor for Code
func NewCode(code string, httpCode int, devMsg, userMsg string) Code {
	return Code{
		code:     code,
		httpCode: httpCode,
		devMsg:   devMsg,
		userMsg:  userMsg,
	}
}

// SetDevMsg set dev msg, return value
func (c *Code) SetDevMsg(msg string) *Code {
	c.devMsg = msg
	return c
}

// SetUserMsg set user msg, return value
func (c *Code) SetUserMsg(msg string) *Code {
	c.userMsg = msg
	return c
}

// DevMsg getter
func (c *Code) DevMsg() string {
	return c.devMsg
}

// UserMsg getter
func (c *Code) UserMsg() string {
	return c.userMsg
}

// Code getter
func (c *Code) Code() string {
	return c.code
}

// HTTPCode getter
func (c *Code) HTTPCode() int {
	return c.httpCode
}

// write any std response code constant here
var (
	// SuccessCode success code
	SuccessCode = Code{
		code:     "SUCCESS",
		httpCode: http.StatusOK,
		devMsg:   "Success",
		userMsg:  "Success",
	}
	// BadRequestErrCode bad request code
	BadRequestErrCode = Code{
		code:     "BAD_REQUEST",
		httpCode: http.StatusBadRequest,
		devMsg:   "Bad Request",
		userMsg:  genericErrorMsg,
	}
	// UnauthorizedCode unauthorized code
	UnauthorizedCode = Code{
		code:     "UNAUTHORIZED",
		httpCode: http.StatusUnauthorized,
		devMsg:   "Unauthorized",
		userMsg:  unauthorizedMsg,
	}
	// PaymentRequiredCode response Code
	PaymentRequiredCode = Code{
		code:     "PAYMENT_REQUIRED",
		httpCode: http.StatusPaymentRequired,
		devMsg:   "Payment Required",
		userMsg:  paymentRequiredMsg,
	}
	// ForbiddenAccessCode forbidden access code
	ForbiddenAccessCode = Code{
		code:     "FORBIDDEN",
		httpCode: http.StatusForbidden,
		devMsg:   "Forbidden",
		userMsg:  forbiddenMsg,
	}
	// NotFoundCode data not found
	NotFoundCode = Code{
		code:     "NOT_FOUND",
		httpCode: http.StatusNotFound,
		devMsg:   "Not Found",
		userMsg:  notFoundMsg,
	}
	// UnprocessableCode unprocessable entity
	UnprocessableCode = Code{
		code:     "UNPROCESSABLE",
		httpCode: http.StatusUnprocessableEntity,
		devMsg:   "Unprocessable Entity",
		userMsg:  unprocessableEntityMsg,
	}
	// TooManyRequestCode rate limit request
	TooManyRequestCode = Code{
		code:     "TOO_MANY_REQUESTS",
		httpCode: http.StatusTooManyRequests,
		devMsg:   "Too Many Requests",
		userMsg:  tooManyRequestMsg,
	}
	// InternalErrCode success code
	InternalErrCode = Code{
		code:     "INTERNAL_SERVER_ERROR",
		httpCode: http.StatusInternalServerError,
		devMsg:   "Internal Server Error",
		userMsg:  genericErrorMsg,
	}
	// NotImplementedCode not implemented code
	NotImplementedCode = Code{
		code:     "NOT_IMPLEMENTED",
		httpCode: http.StatusNotImplemented,
		devMsg:   "Not Implemented",
		userMsg:  "Not Implemented",
	}
)
