package errs

import "net/http"

func init() {
	setDefaults()
}

var (
	genericUserMsg = "Uh oh! Something unexpected occurred. Please send a report to support@your-app.com so we can fix it."
)

var (
	ErrTypeExpected *ErrType
	ErrTypeHTTP     *ErrType
)

var (
	// 4xx
	ErrTypeBadRequest *StatusType
	NewErrBadRequest  func() *Status

	ErrTypeUnauthorized *StatusType
	NewErrUnauthorized  func() *Status

	ErrTypePaymentRequired *StatusType
	NewErrPaymentRequired  func() *Status

	ErrTypeForbiddenAccess *StatusType
	NewErrForbiddenAccess  func() *Status

	ErrTypeNotFound *StatusType
	NewErrNotFound  func() *Status

	ErrTypeUnprocessable *StatusType
	NewErrUnprocessable  func() *Status

	ErrTypeTooManyRequest *StatusType
	NewErrTooManyRequest  func() *Status

	// 5xx
	ErrTypeInternalErr *StatusType
	NewErrInternalErr  func() *Status

	ErrTypeNotImplemented *StatusType
	NewErrNotImplemented  func() *Status
)

func setDefaults() {
	ErrTypeExpected = New("expected error").
		WithSeverity(SeverityWarning).
		WithUserMsg("We're experiencing a hiccup in our system. Please try again in a minute."). // TODO: better msg
		Freeze()

	ErrTypeHTTP = New("http error").Freeze()

	// 4xx
	ErrTypeBadRequest = NewStatus("BAD_REQUEST",
		New("bad request").
			WithSeverity(SeverityError).
			WithPublicMsg("Bad Request").
			WithUserMsg(genericUserMsg).
			AddBase(ErrTypeHTTP)).
		WithHTTP(http.StatusBadRequest).
		Freeze()
	NewErrBadRequest = ErrTypeBadRequest.New

	ErrTypeUnauthorized = NewStatus("UNAUTHORIZED",
		New("unauthorized").
			WithSeverity(SeverityError).
			WithPublicMsg("Unauthorized").
			WithUserMsg("Sorry, you need to be logged in to access this page. Please log in and try again.").
			AddBase(ErrTypeHTTP)).
		WithHTTP(http.StatusUnauthorized).
		Freeze()
	NewErrUnauthorized = ErrTypeUnauthorized.New

	ErrTypePaymentRequired = NewStatus("PAYMENT_REQUIRED",
		New("payment required").
			WithSeverity(SeverityWarning).
			WithPublicMsg("Payment Required").
			WithUserMsg("You need to upgrade your plan to access this feature.").
			AddBase(ErrTypeHTTP, ErrTypeExpected)).
		WithHTTP(http.StatusPaymentRequired).
		Freeze()
	NewErrPaymentRequired = ErrTypePaymentRequired.New

	ErrTypeForbiddenAccess = NewStatus("FORBIDDEN",
		New("forbidden").
			WithSeverity(SeverityError).
			WithPublicMsg("Forbidden").
			WithUserMsg("You do not have the necessary permissions to view this item.").
			AddBase(ErrTypeHTTP)).
		WithHTTP(http.StatusForbidden).
		Freeze()
	NewErrForbiddenAccess = ErrTypeForbiddenAccess.New

	ErrTypeNotFound = NewStatus("NOT_FOUND",
		New("not found").
			WithSeverity(SeverityError).
			WithPublicMsg("Not Found").
			WithUserMsg("Oops! We couldn't find what you were looking for. Please send a report to support@your-app.com if you believe this was an error.").
			AddBase(ErrTypeHTTP)).
		WithHTTP(http.StatusNotFound).
		Freeze()
	NewErrNotFound = ErrTypeNotFound.New

	ErrTypeUnprocessable = NewStatus("UNPROCESSABLE",
		New("unprocessable entity").
			WithSeverity(SeverityError).
			WithPublicMsg("Unprocessable Entity").
			WithUserMsg("The request could not be processed correctly due to a mistake in the information provided. Please review and try again.").
			AddBase(ErrTypeHTTP)).
		WithHTTP(http.StatusUnprocessableEntity).
		Freeze()
	NewErrUnprocessable = ErrTypeUnprocessable.New

	ErrTypeTooManyRequest = NewStatus("TOO_MANY_REQUESTS",
		New("too many requests").
			WithSeverity(SeverityWarning).
			WithPublicMsg("Too Many Requests").
			WithUserMsg("You've made too many requests. Please take a break and try again later.").
			AddBase(ErrTypeHTTP, ErrTypeExpected)).
		WithHTTP(http.StatusTooManyRequests).
		Freeze()
	NewErrTooManyRequest = ErrTypeTooManyRequest.New

	// 5xx
	ErrTypeInternalErr = NewStatus("INTERNAL_SERVER_ERROR",
		New("internal server error").
			WithSeverity(SeverityError).
			WithPublicMsg("Internal Server Error").
			WithUserMsg(genericUserMsg).
			AddBase(ErrTypeHTTP)).
		WithHTTP(http.StatusInternalServerError).
		Freeze()
	NewErrInternalErr = ErrTypeInternalErr.New

	ErrTypeNotImplemented = NewStatus("NOT_IMPLEMENTED",
		New("not implemented").
			WithSeverity(SeverityWarning).
			WithPublicMsg("Not Implemented").
			WithUserMsg(genericUserMsg).
			AddBase(ErrTypeHTTP, ErrTypeExpected)).
		WithHTTP(http.StatusNotImplemented).
		Freeze()
	NewErrNotImplemented = ErrTypeNotImplemented.New
}
