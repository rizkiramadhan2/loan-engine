package errs

import (
	"errors"
	"fmt"
	"testing"
)

func TestErr_Error(t *testing.T) {
	type fields struct {
		o Options
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "no cause, no internal msg - empty string",
			fields: fields{
				o: Options{
					Cause:   nil,
					Context: nil,

					Severity: SeverityUnknown,

					InternalMsg: "",
					PublicMsg:   "",
					UserMsg:     "",

					StackTrace: nil,
				},
			},
			want: "",
		},
		{
			name: "some error, no internal msg - error string",
			fields: fields{
				o: Options{
					Cause:   New("some error"),
					Context: nil,

					Severity: SeverityUnknown,

					InternalMsg: "",
					PublicMsg:   "",
					UserMsg:     "",

					StackTrace: nil,
				},
			},
			want: "some error",
		},
		{
			name: "no cause, some internal msg - internal msg",
			fields: fields{
				o: Options{
					Cause:   nil,
					Context: nil,

					Severity: SeverityUnknown,

					InternalMsg: "some internal msg",
					PublicMsg:   "",
					UserMsg:     "",

					StackTrace: nil,
				},
			},
			want: "some internal msg",
		},
		{
			name: "some error, some internal msg - internal msg",
			fields: fields{
				o: Options{
					Cause:   New("some error"),
					Context: nil,

					Severity: SeverityUnknown,

					InternalMsg: "some internal msg",
					PublicMsg:   "",
					UserMsg:     "",
				},
			},
			want: "some internal msg: some error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrType{
				o: &tt.fields.o,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Err.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErr_Unwrap(t *testing.T) {
	type fields struct {
		o Options
	}
	err1 := errors.New("some error")
	err2 := fmt.Errorf("some error: %w", err1)
	err3 := Wrap(err2, "some internal msg")
	tests := []struct {
		name     string
		fields   fields
		wantErrs []error
	}{
		{
			name: "1 wrap",
			fields: fields{
				o: Options{
					Cause:   err1,
					Context: nil,

					Severity: SeverityUnknown,

					InternalMsg: "",
					PublicMsg:   "",
					UserMsg:     "",

					StackTrace: nil,
				},
			},
			wantErrs: []error{err1},
		},
		{
			name: "2 wraps",
			fields: fields{
				o: Options{
					Cause:   err2,
					Context: nil,

					Severity: SeverityUnknown,

					InternalMsg: "",
					PublicMsg:   "",
					UserMsg:     "",

					StackTrace: nil,
				},
			},
			wantErrs: []error{err1, err2},
		},
		{
			name: "3 wraps",
			fields: fields{
				o: Options{
					Cause:   err3,
					Context: nil,

					Severity: SeverityUnknown,

					InternalMsg: "",
					PublicMsg:   "",
					UserMsg:     "",

					StackTrace: nil,
				},
			},
			wantErrs: []error{err1, err2, err3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrType{
				o: &tt.fields.o,
			}
			for i, err := range tt.wantErrs {
				if !errors.Is(e, err) {
					t.Errorf("Err.Unwrap() = %v, want %v, i=%v", e, err, i)
				}
			}
		})
	}
}

func TestWrap(t *testing.T) {
	err1 := errors.New("some error")

	e := Wrap(err1, "some error 2")

	if !errors.Is(e, err1) {
		t.Errorf("Wrap() = %v, want %v", e, err1)
	}

	if e.Error() != "some error 2: some error" {
		t.Errorf("Wrap.Error() = %v, want %v", e.Error(), "some error 2: some error")
	}
}

func TestErr_Getters(t *testing.T) {
	type fields struct {
		o    Options
		data int
	}
	err1 := errors.New("some error")
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "all fields",
			fields: fields{
				o: Options{
					Cause:   err1,
					Context: nil,

					Severity: SeverityUnknown,

					InternalMsg: "some internal msg",
					PublicMsg:   "some public msg",
					UserMsg:     "some user msg",
				},
				data: 420,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrType{
				o:    &tt.fields.o,
				data: tt.fields.data,
			}
			if got := e.Cause(); !errors.Is(got, err1) {
				t.Errorf("Err.Cause() = %v, want %v", got, err1)
			}
			if got := e.Context(); got != nil {
				t.Errorf("Err.Context() = %v, want %v", got, nil)
			}
			if got := e.Severity(); got != SeverityUnknown {
				t.Errorf("Err.Severity() = %v, want %v", got, SeverityUnknown)
			}
			if got := e.InternalMsg(); got != "some internal msg" {
				t.Errorf("Err.InternalMsg() = %v, want %v", got, "some internal msg: some error")
			}
			if got := e.PublicMsg(); got != "some public msg" {
				t.Errorf("Err.PublicMsg() = %v, want %v", got, "some public msg")
			}
			if got := e.UserMsg(); got != "some user msg" {
				t.Errorf("Err.UserMsg() = %v, want %v", got, "some user msg")
			}
			if got := e.Data(); got != 420 {
				t.Errorf("Err.Data() = %v, want %v", got, 420)
			}
		})
	}
}
