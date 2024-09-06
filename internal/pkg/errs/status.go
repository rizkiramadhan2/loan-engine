package errs

import (
	"errors"

	"simple-app/internal/pkg/syncg"
)

var (
	mapStringStatus = syncg.Map[string, *Status]{}
)

func GetStatus(code string) (status *Status, ok bool) {
	status, ok = mapStringStatus.Load(code)
	return
}

type StatusData struct {
	Code    string
	isValid bool

	httpStatusCode int
	grpcStatusCode int
}

type Status struct {
	*Err
	StatusData
}

type StatusType struct {
	*ErrType
	StatusData
}

type HTTPStatus struct {
	Status
}

func (s *HTTPStatus) HTTPStatusCode() (code int) {
	return s.httpStatusCode
}

type GRPCStatus struct {
	Status
}

func (s *GRPCStatus) GRPCStatusCode() (code int) {
	return s.grpcStatusCode
}

func (s *StatusType) Unwrap() error {
	return s.ErrType
}

func (s *Status) Unwrap() error {
	return s.Err
}

func (s *Status) As(target any) bool {
	switch t := target.(type) {
	case **HTTPStatus:
		if !s.isValid || s.httpStatusCode == -1 {
			return false
		}
		*t = &HTTPStatus{
			Status: *s,
		}
		return true
	case **GRPCStatus:
		if !s.isValid || s.grpcStatusCode == -1 {
			return false
		}
		*t = &GRPCStatus{
			Status: *s,
		}
	}
	return false
}

func (s *Status) Is(target error) bool {
	if target, ok := target.(*Status); ok {
		if s.Code == target.Code {
			return true
		}
	}

	if target, ok := target.(*StatusType); ok {
		if s.Code == target.Code {
			return true
		}
	}

	if errors.Is(s.Err, target) {
		return true
	}

	return false
}

func (s *Status) save() {
	mapStringStatus.Store(s.Code, s)
}

func NewStatus(code string, err *Err) *Status {
	s := Status{
		Err: err,
		StatusData: StatusData{
			Code:    code,
			isValid: true,

			httpStatusCode: -1,
			grpcStatusCode: -1,
		},
	}
	defer s.save()

	return &s
}

func (s *StatusType) Copy() *Status {
	return &Status{
		Err:        s.ErrType.Copy(),
		StatusData: s.StatusData,
	}
}

func (s *StatusType) New() *Status {
	as := &Status{
		Err:        NewWithBase(s, s.ErrType.InternalMsg()).WithData(s.data),
		StatusData: s.StatusData,
	}
	return as
}

func (s *Status) Freeze() *StatusType {
	return &StatusType{
		ErrType:    s.Err.Freeze(),
		StatusData: s.StatusData,
	}
}

func (s *Status) WithHTTP(code int) *Status {
	defer s.save()
	s.httpStatusCode = code
	return s
}

func (s *Status) WithGRPC(code int) *Status {
	defer s.save()
	s.grpcStatusCode = code
	return s
}
