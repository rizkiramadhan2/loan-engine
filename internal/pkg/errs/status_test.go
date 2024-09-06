package errs

import (
	"errors"
	"testing"
)

func TestNewStatus(t *testing.T) {
	randomHTTPErr := NewStatus("status", New("http err")).WithHTTP(1)

	wrappedHTTP := Wrap(randomHTTPErr, "wrap")
	if !errors.Is(wrappedHTTP, randomHTTPErr) {
		t.Error("wrappedHTTP is not randomHTTPErr")
	}

	randomErr := errors.New("random err")
	wrappedHTTP = wrappedHTTP.WithCause(randomErr)
	if !errors.Is(wrappedHTTP, randomErr) {
		t.Error("wrapped is not randomErr")
	}

	status := &Status{}
	if ok := errors.As(NewErrBadRequest(), &status); !ok {
		t.Error("ErrBadRequest() is not of type Status")
	}

	httpErr := &HTTPStatus{}
	if ok := errors.As(NewErrBadRequest(), &httpErr); !ok {
		t.Error("ErrBadRequest() is not of type HTTPStatus")
	}

	if ok := errors.Is(NewErrBadRequest(), ErrTypeBadRequest); !ok {
		t.Error("ErrBadRequest() is not ErrTypeBadRequest")
	}

	if ok := errors.Is(NewErrBadRequest().WithCause(New("err")), ErrTypeBadRequest); !ok {
		t.Error("ErrBadRequest() with cause is not ErrTypeBadRequest")
	}

	if ok := errors.Is(NewErrBadRequest().WithCause(errors.New("err")), ErrTypeBadRequest); !ok {
		t.Error("ErrBadRequest() with cause is not ErrTypeBadRequest")
	}

	grpcErr := &GRPCStatus{}
	if ok := errors.As(NewErrPaymentRequired(), &grpcErr); ok {
		t.Error("ErrPaymentRequired() is of type GRPCStatus")
	}
}

func TestGetStatus(t *testing.T) {
	if status, ok := GetStatus(ErrTypeBadRequest.Code); !ok {
		t.Error("ErrTypeBadRequest.Code is not found")
	} else {
		if !errors.Is(status, ErrTypeBadRequest) {
			t.Error("status is not ErrTypeBadRequest")
		}
	}

	if _, ok := GetStatus("random"); ok {
		t.Error("random is found")
	}
}
