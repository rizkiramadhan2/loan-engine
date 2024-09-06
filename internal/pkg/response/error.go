package response

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
)

// interface implementation check
var (
	_ error                       = &Error{}
	_ interface{ Unwrap() error } = &Error{}
	_ interface{ Is(error) bool } = &Error{}
)

// Error struct
type Error struct {
	// for response
	code    Code
	userMsg string
	data    interface{}

	// for dev
	cause      error
	errMsg     string
	stackTrace []string
}

// ErrIs compares error code
func ErrIs(err error, code Code) bool {
	if e, ok := err.(*Error); ok {
		return e.code.code == code.code
	}
	return false
}

// UnwrapErr unwraps error
func UnwrapErr(err error) error {
	if e, ok := err.(*Error); ok {
		return e.cause
	}
	return err
}

// Is implementation
func (e *Error) Is(target error) bool {
	if target, ok := target.(*Error); ok {
		if e.code.code == target.code.code {
			return true
		}
	}

	if errors.Is(e.cause, target) {
		return true
	}

	return false
}

// Code getter
func (e *Error) Code() Code {
	return e.code
}

// Cause unwraps the error
func (e *Error) Cause() error {
	return e.cause
}

// Unwrap unwraps the error
func (e *Error) Unwrap() error {
	return e.cause
}

// Error implementation of error interface
func (e Error) Error() string {
	if e.errMsg == "" {
		return e.cause.Error()
	}
	return e.errMsg + ": " + e.cause.Error()
}

// Readable return readable error msg
func (e Error) Readable() string {
	return e.userMsg
}

// StackTrace return stack trace
func (e *Error) StackTrace() []string {
	return e.stackTrace
}

// SetUserMsg sets the message
func (e *Error) SetUserMsg(msg string) *Error {
	e.userMsg = msg
	return e
}

// WrapUserMsg sets the message
func (e *Error) WrapUserMsg(msg string) *Error {
	e.userMsg = msg + ": " + e.userMsg
	return e
}

// WithData sets the message
func (e *Error) WithData(data interface{}) *Error {
	e.data = data
	return e
}

func wrapErrCode(err error, code Code, msg ...string) *Error {
	newErr, ok := err.(*Error)
	if !ok {
		newErr = &Error{
			code:    code,
			userMsg: code.userMsg,

			cause:      err,
			errMsg:     "",
			stackTrace: makeTrace(1),
		}
	}

	// replace msg from current stack
	if len(msg) > 0 {
		if newErr.errMsg == "" {
			newErr.errMsg = msg[0]
		} else {
			newErr.errMsg = msg[0] + ": " + newErr.errMsg
		}
	}

	return newErr
}

// NewError create std response.Error
func NewError(cause string, code ...Code) *Error {
	err := errors.New(cause)
	c := InternalErrCode
	if len(code) > 0 {
		c = code[0]
	}

	return wrapErrCode(err, c)
}

func makeErr(code Code, err error) *Error {
	if err == nil {
		return NewError(code.DevMsg(), code)
	}
	return WrapErrCode(err, code, code.DevMsg())
}

// NewErrorBuilder create std response.Error
func NewErrorBuilder(code Code) func(err error) *Error {
	return func(err error) *Error {
		return makeErr(code, err)
	}
}

// WrapErrCode wrap error and tampering current response code
func WrapErrCode(err error, code Code, msg ...string) *Error {
	return wrapErrCode(err, code, msg...)
}

// WrapErr wrap error without tampering existing response code if any
func WrapErr(err error, msg ...string) *Error {
	return wrapErrCode(err, InternalErrCode, msg...)
}

// DeferWrap wrap in place
func DeferWrap(err *error, msg ...string) {
	if err == nil || *err == nil {
		return
	}

	*err = wrapErrCode(*err, InternalErrCode, msg...)
}

func makeTrace(skip int) []string {
	trace := []string{}
	pc := make([]uintptr, 15)
	n := runtime.Callers(skip+2, pc)
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, next := frames.Next()
		t := fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function)
		trace = append(trace, t)
		if !next {
			break
		}
	}
	return trace
}

var (
	InvalidRequestPayloadCode = func(val ...string) *Error {
		if len(val) == 0 {
			val = []string{"invalid payload format"}
		}

		return NewError(val[0], BadRequestErrCode).WithData(val)
	}
)
