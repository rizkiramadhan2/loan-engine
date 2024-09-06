package errs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type Opts struct {
	GenericErrorMsg string
}

func Init(opts Opts) {
	if opts.GenericErrorMsg != "" {
		genericUserMsg = opts.GenericErrorMsg
	}
	setDefaults()
}

type Options struct {
	BaseTypes []error
	Cause     error
	Context   context.Context

	Severity Severity

	InternalMsg string
	PublicMsg   string
	UserMsg     string

	StackTrace []string
}

type ErrType struct {
	o    *Options
	data any
}

type Err struct {
	ErrType
}

// interface implementation check
var (
	_ error                       = New("")
	_ interface{ Unwrap() error } = New("")
	_ interface{ Is(error) bool } = New("")
)

// New constructor for Error
func New(internalMsg string) *Err {
	return NewWithOpts(Options{
		Cause:   nil,
		Context: context.Background(),

		Severity: SeverityUnknown,

		InternalMsg: internalMsg,
		PublicMsg:   "an error occurred",
		UserMsg:     genericUserMsg,

		StackTrace: makeTrace(1),
	})
}

// NewWithOpts constructor for Error
func NewWithOpts(o Options) *Err {
	return (&Err{
		ErrType: ErrType{
			o: &o,
		},
	}).WithStackTrace()
}

// Error implementation
func (e *ErrType) Error() string {
	msgs := []string{}

	if e.o.InternalMsg != "" {
		msgs = append(msgs, e.o.InternalMsg)
	}

	if e.o.Cause != nil {
		em := e.o.Cause.Error()
		msgs = append(msgs, em)
	}

	if len(msgs) == 0 && len(e.o.BaseTypes) > 0 {
		if len(e.o.BaseTypes) == 1 {
			msgs = append(msgs, e.o.BaseTypes[0].Error())
		} else {
			bases := []string{}
			for i := range e.o.BaseTypes {
				b := e.o.BaseTypes[i]
				bases = append(bases, b.Error())
			}
			slices.Sort(bases)
			msgs = append(msgs, strings.Join(bases, ","))
		}
	}

	msg := strings.Join(msgs, ": ")

	return msg

}

// Unwrap implementation
func (e *ErrType) Unwrap() error {
	return e.o.Cause
}

// Is implementation
func (e *ErrType) Is(target error) bool {
	if target, ok := target.(*Err); ok {
		if e.o == target.o {
			return true
		}
	}

	if errors.Is(e.o.Cause, target) {
		return true
	}

	for _, b := range e.o.BaseTypes {
		if errors.Is(b, target) {
			return true
		}
	}

	return false
}

// As implementation
func (e *Err) As(target any) bool {
	if target, ok := target.(**Err); ok {
		*target = e
		return true
	}
	if target, ok := target.(**ErrType); ok {
		*target = &e.ErrType
		return true
	}
	return false
}

// As implementation
func (e *ErrType) As(target any) bool {
	if target, ok := target.(**ErrType); ok {
		*target = e
		return true
	}
	if target, ok := target.(**Err); ok {
		*target = e.Copy()
		return true
	}
	return false
}

// Copy copies an error
func (e *ErrType) Copy() *Err {
	return NewWithOpts(*e.o)
}

// CopyAsBase copies an error with current error as base
func (e *ErrType) CopyAsBase() *Err {
	return NewWithBase(e, "").WithData(e.data)
}

// Wrap wraps an error
func Wrap(cause error, internalMsg string) *Err {
	if err, ok := cause.(*Err); ok {
		return NewWithOpts(*err.o).WithCause(cause).WithInternalMsg(internalMsg)
	}
	return New(internalMsg).WithCause(cause)
}

// DeferWrap wrap in place called like `defer errs.DeferWrap(&err)` at the start of function
// function must have named return parameter
func DeferWrap(err *error, msg ...string) {
	if err == nil || *err == nil {
		return
	}

	m := ""
	if len(msg) > 0 {
		m = msg[0]
	}

	*err = Wrap(*err, m)
}

// NewWithBase wraps an error
func NewWithBase(base error, internalMsg string) *Err {
	var e *ErrType
	if ok := errors.As(base, &e); ok {
		return NewWithOpts(*e.o).WithBase([]error{base}).WithInternalMsg(internalMsg)
	}
	return New(internalMsg).WithBase([]error{base})
}

// Freeze freezes the error
func (e *Err) Freeze() *ErrType {
	return &e.ErrType
}

// AddBase sets the base of the error
func (e *Err) AddBase(base ...error) *Err {
	e.o.BaseTypes = append(e.o.BaseTypes, base...)
	return e
}

// WithBase sets the base of the error
func (e *Err) WithBase(bases []error) *Err {
	e.o.BaseTypes = bases
	return e
}

// WithCause sets the cause of the error
func (e *Err) WithCause(cause error) *Err {
	e.o.Cause = cause
	if ep, ok := cause.(*Err); ok {
		e.o.StackTrace = ep.o.StackTrace
	}
	return e
}

// WithContext sets the context of the error
func (e *Err) WithContext(ctx context.Context) *Err {
	e.o.Context = ctx
	return e
}

// WithSeverity sets the severity of the error
func (e *Err) WithSeverity(severity Severity) *Err {
	e.o.Severity = severity
	return e
}

// WithInternalMsg sets the internal message of the error
func (e *Err) WithInternalMsg(internalMsg string) *Err {
	e.o.InternalMsg = internalMsg
	return e
}

// WithInternalMsgf sets the internal message of the error
func (e *Err) WithInternalMsgf(format string, args ...any) *Err {
	e.o.InternalMsg = fmt.Sprintf(format, args...)
	return e
}

// WithPublicMsg sets the API message of the error
func (e *Err) WithPublicMsg(apiMsg string) *Err {
	e.o.PublicMsg = apiMsg
	return e
}

// WithUserMsg sets the user message of the error
func (e *Err) WithUserMsg(userMsg string) *Err {
	e.o.UserMsg = userMsg
	return e
}

// WithStackTrace sets the stack trace of the error
func (e *Err) WithStackTrace() *Err {
	e.o.StackTrace = makeTrace(1)
	return e
}

// WithData sets the data of the error
func (e *Err) WithData(data any) *Err {
	e.data = data
	return e
}

// Cause getter
func (e *ErrType) Cause() error {
	return e.o.Cause
}

// Context getter
func (e *ErrType) Context() context.Context {
	return e.o.Context
}

// Severity getter
func (e *ErrType) Severity() Severity {
	return e.o.Severity
}

// InternalMsg getter
func (e *ErrType) InternalMsg() string {
	return e.o.InternalMsg
}

// APIMsg getter
func (e *ErrType) PublicMsg() string {
	return e.o.PublicMsg
}

// UserMsg getter
func (e *ErrType) UserMsg() string {
	return e.o.UserMsg
}

// StackTrace getter
func (e *ErrType) StackTrace() []string {
	return e.o.StackTrace
}

// Data getter
func (e *ErrType) Data() any {
	return e.data
}
