package errs

import (
	"errors"
	"testing"
)

func TestHTTPFields(t *testing.T) {
	statusErr := NewErrBadRequest()
	var err error = statusErr
	want := "bad request"
	if err.Error() != want {
		t.Errorf("err.Error() != \"%s\" got \"%s\"", want, err.Error())
	}

	errr := errors.New("random err")
	err = NewErrBadRequest().WithCause(errr)
	want = "bad request: random err"
	if err.Error() != want {
		t.Errorf("err.Error() != \"%s\" got \"%s\"", want, err.Error())
	}

	want = NewErrBadRequest().UserMsg()
	if statusErr.UserMsg() != want {
		t.Errorf("status.UserMsg() != \"%s\" got \"%s\"", want, statusErr.UserMsg())
	}

	want = "BAD_REQUEST"
	var status *HTTPStatus
	if ok := errors.As(statusErr, &status); !ok {
		t.Error("statusErr is not of type HTTPStatus")
	} else {
		if status.Code != want {
			t.Errorf("status.Code() != \"%s\" got \"%s\"", want, status.Code)
		}
	}
}

func TestExpected(t *testing.T) {
	err := NewErrPaymentRequired()
	if !errors.Is(err, ErrTypeExpected) {
		t.Errorf("wrap expected doesn't work")
	}
}
