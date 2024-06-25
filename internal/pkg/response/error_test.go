package response

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapErrCode(t *testing.T) {
	type args struct {
		err  error
		code Code
		msg  []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantMsg string
	}{
		{
			name: "Test Wrap w/o Readable Msg",
			args: args{
				err:  errors.New("some err"),
				code: InternalErrCode,
				msg:  nil,
			},
			wantErr: true,
			wantMsg: "some err",
		}, {
			name: "Test Wrap",
			args: args{
				err:  WrapErr(errors.New("some err")),
				code: NotFoundCode,
				msg:  nil,
			},
			wantMsg: "some err",
			wantErr: true,
		}, {
			name: "Test Wrap 2x",
			args: args{
				err:  WrapErr(errors.New("some err"), "test"),
				code: NotFoundCode,
				msg:  []string{"wrap"},
			},
			wantMsg: "wrap: test: some err",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := WrapErrCode(tt.args.err, tt.args.code, tt.args.msg...)
			if (err != nil) != tt.wantErr {
				t.Errorf("WrapErr() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.wantMsg, err.Error())
		})
	}
}

func TestErrorFunc(t *testing.T) {
	type fields struct {
		code       Code
		msg        string
		err        error
		stackTrace []string
	}
	tests := []struct {
		name           string
		fields         fields
		want           string
		wantReadable   string
		wantStackTrace []string
	}{
		{
			name: "Test 1",
			fields: fields{
				code:       Code{},
				err:        errors.New("some err"),
				stackTrace: []string{},
				msg:        "this is readable message",
			},
			want:           "some err",
			wantReadable:   "this is readable message",
			wantStackTrace: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				code:    tt.fields.code,
				userMsg: tt.fields.msg,

				cause:      tt.fields.err,
				errMsg:     tt.fields.err.Error(),
				stackTrace: tt.fields.stackTrace,
			}

			assert.Equal(t, tt.want, e.Error())
			assert.Equal(t, tt.wantReadable, e.Readable())
			assert.Equal(t, tt.wantStackTrace, e.StackTrace())
		})
	}
}

func TestDeferWrap(t *testing.T) {
	type args struct {
		err func() *error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test Nil Err",
			args: args{
				err: func() *error {
					return nil
				},
			},
		}, {
			name: "Test Std Error",
			args: args{
				err: func() *error {
					e := errors.New("some err")
					return &e
				},
			},
		}, {
			name: "Test Response Error",
			args: args{
				err: func() *error {
					var e error = WrapErr(errors.New("some err"))
					return &e
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DeferWrap(tt.args.err())
		})
	}
}

func TestErrIs(t *testing.T) {
	type args struct {
		err  error
		code Code
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				err:  NewError("error", BadRequestErrCode),
				code: BadRequestErrCode,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrIs(tt.args.err, tt.args.code); got != tt.want {
				t.Errorf("ErrIs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Cause(t *testing.T) {
	type fields struct {
		code       Code
		userMsg    string
		data       interface{}
		cause      error
		errMsg     string
		stackTrace []string
	}
	err1 := errors.New("err")
	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name: "success",
			fields: fields{
				cause: err1,
			},
			want: err1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				code:       tt.fields.code,
				userMsg:    tt.fields.userMsg,
				data:       tt.fields.data,
				cause:      tt.fields.cause,
				errMsg:     tt.fields.errMsg,
				stackTrace: tt.fields.stackTrace,
			}
			if err := e.Cause(); err != tt.want {
				t.Errorf("Error.Cause() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		code       Code
		userMsg    string
		data       interface{}
		cause      error
		errMsg     string
		stackTrace []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success",
			fields: fields{
				errMsg: "error",
			},
			want: "error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				code:       tt.fields.code,
				userMsg:    tt.fields.userMsg,
				data:       tt.fields.data,
				cause:      tt.fields.cause,
				errMsg:     tt.fields.errMsg,
				stackTrace: tt.fields.stackTrace,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Readable(t *testing.T) {
	type fields struct {
		code       Code
		userMsg    string
		data       interface{}
		cause      error
		errMsg     string
		stackTrace []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success",
			fields: fields{
				userMsg: "msg",
			},
			want: "msg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				code:       tt.fields.code,
				userMsg:    tt.fields.userMsg,
				data:       tt.fields.data,
				cause:      tt.fields.cause,
				errMsg:     tt.fields.errMsg,
				stackTrace: tt.fields.stackTrace,
			}
			if got := e.Readable(); got != tt.want {
				t.Errorf("Error.Readable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_StackTrace(t *testing.T) {
	type fields struct {
		code       Code
		userMsg    string
		data       interface{}
		cause      error
		errMsg     string
		stackTrace []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "success",
			fields: fields{
				stackTrace: []string{"1", "2"},
			},
			want: []string{"1", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				code:       tt.fields.code,
				userMsg:    tt.fields.userMsg,
				data:       tt.fields.data,
				cause:      tt.fields.cause,
				errMsg:     tt.fields.errMsg,
				stackTrace: tt.fields.stackTrace,
			}
			if got := e.StackTrace(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.StackTrace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_SetUserMsg(t *testing.T) {
	type fields struct {
		code       Code
		userMsg    string
		data       interface{}
		cause      error
		errMsg     string
		stackTrace []string
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Error
	}{
		{
			name:   "success",
			fields: fields{},
			args: args{
				msg: "msg1",
			},
			want: &Error{
				userMsg: "msg1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				code:       tt.fields.code,
				userMsg:    tt.fields.userMsg,
				data:       tt.fields.data,
				cause:      tt.fields.cause,
				errMsg:     tt.fields.errMsg,
				stackTrace: tt.fields.stackTrace,
			}
			if got := e.SetUserMsg(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.SetUserMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_WrapUserMsg(t *testing.T) {
	type fields struct {
		code       Code
		userMsg    string
		data       interface{}
		cause      error
		errMsg     string
		stackTrace []string
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Error
	}{
		{
			name: "success no data",
			fields: fields{
				userMsg: "msg1",
			},
			args: args{
				msg: "msg2",
			},
			want: &Error{
				userMsg: "msg2: msg1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				code:       tt.fields.code,
				userMsg:    tt.fields.userMsg,
				data:       tt.fields.data,
				cause:      tt.fields.cause,
				errMsg:     tt.fields.errMsg,
				stackTrace: tt.fields.stackTrace,
			}
			if got := e.WrapUserMsg(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.WrapUserMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_WithData(t *testing.T) {
	type fields struct {
		code       Code
		userMsg    string
		data       interface{}
		cause      error
		errMsg     string
		stackTrace []string
	}
	type args struct {
		data interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Error
	}{
		{
			name:   "success",
			fields: fields{},
			args: args{
				data: "data",
			},
			want: &Error{
				data: "data",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				code:       tt.fields.code,
				userMsg:    tt.fields.userMsg,
				data:       tt.fields.data,
				cause:      tt.fields.cause,
				errMsg:     tt.fields.errMsg,
				stackTrace: tt.fields.stackTrace,
			}
			if got := e.WithData(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.WithData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wrapErrCode(t *testing.T) {
	type args struct {
		err  error
		code Code
		msg  []string
	}
	err1 := errors.New("error")
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "success",
			args: args{
				err:  err1,
				code: BadRequestErrCode,
				msg:  []string{"msg"},
			},
			want: &Error{
				code:       BadRequestErrCode,
				userMsg:    BadRequestErrCode.userMsg,
				data:       nil,
				cause:      err1,
				errMsg:     "msg: " + err1.Error(),
				stackTrace: nil, //TODO:
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wrapErrCode(tt.args.err, tt.args.code, tt.args.msg...)
			got.stackTrace = nil
			// got.cause = nil
			// tt.want.cause = nil
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wrapErrCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewError(t *testing.T) {
	type args struct {
		msg  string
		code []Code
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "success",
			args: args{
				msg:  "msg",
				code: []Code{BadRequestErrCode},
			},
			want: &Error{
				code:       BadRequestErrCode,
				userMsg:    BadRequestErrCode.userMsg,
				data:       nil,
				cause:      nil,
				errMsg:     "msg",
				stackTrace: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewError(tt.args.msg, tt.args.code...)
			got.stackTrace = nil
			got.cause = nil
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrapErr(t *testing.T) {
	type args struct {
		err error
		msg []string
	}
	err1 := errors.New("err")
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "success",
			args: args{
				err: err1,
				msg: []string{"msg"},
			},
			want: &Error{
				code:       InternalErrCode,
				userMsg:    InternalErrCode.userMsg,
				data:       nil,
				cause:      err1,
				errMsg:     "msg: err",
				stackTrace: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WrapErr(tt.args.err, tt.args.msg...)
			got.stackTrace = nil
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WrapErr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeTrace(t *testing.T) {
	type args struct {
		skip int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeTrace(tt.args.skip); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeTrace() = %v, want %v", got, tt.want)
			}
		})
	}
}
