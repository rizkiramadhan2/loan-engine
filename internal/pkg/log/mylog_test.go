package log

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"

	goerrors "github.com/go-errors/errors"
)

func TestErrorfCtx(t *testing.T) {
	type args struct {
		ctx    context.Context
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test Error Ctx",
			args: args{
				ctx:    context.Background(),
				format: "%s",
				args: []interface{}{
					"test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ErrorfCtx(tt.args.ctx, tt.args.format, tt.args.args...)
		})
		t.Run(tt.name, func(t *testing.T) {
			InfofCtx(tt.args.ctx, tt.args.format, tt.args.args...)
		})
		t.Run(tt.name, func(t *testing.T) {
			ErrorfCtxWithSkip(tt.args.ctx, 1, tt.args.format, tt.args.args...)
		})
		t.Run(tt.name, func(t *testing.T) {
			InfofCtxWithSkip(tt.args.ctx, 1, tt.args.format, tt.args.args...)
		})
	}
}

func TestHandleDefer(t *testing.T) {
	type args struct {
		ctx     context.Context
		logInfo *LogInfoType
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test Defer Log",
			args: args{
				ctx:     context.Background(),
				logInfo: &LogInfoType{},
			},
		}, {
			name: "Test Defer Log with Field",
			args: args{
				ctx: context.Background(),
				logInfo: &LogInfoType{
					Fields:   map[string]interface{}{"test": "test"},
					IsLogged: true,
					Err:      goerrors.New("some err"),
				},
			},
		}, {
			name: "Test Defer Log with Field",
			args: args{
				ctx: context.Background(),
				logInfo: &LogInfoType{
					Fields:   map[string]interface{}{"test": "test"},
					IsLogged: true,
					Err:      errors.New("some err"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleDefer(tt.args.ctx, tt.args.logInfo)
		})
	}
}

func TestLogInfoType_GetLatency(t *testing.T) {
	now := time.Now().Add(-2 * time.Millisecond)
	type fields struct {
		Name        string
		Request     interface{}
		Response    interface{}
		Err         error
		Latency     string
		requestTime time.Time
		Fields      Fields
		ErrorCode   int
		RequestID   string
		IsLogged    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test Success",
			fields: fields{
				Name:        "Some Log",
				Request:     nil,
				Response:    nil,
				Err:         nil,
				Latency:     "",
				requestTime: now,
				Fields:      nil,
				ErrorCode:   0,
				RequestID:   "",
				IsLogged:    false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LogInfoType{
				Name:        tt.fields.Name,
				Request:     tt.fields.Request,
				Response:    tt.fields.Response,
				Err:         tt.fields.Err,
				Latency:     tt.fields.Latency,
				requestTime: tt.fields.requestTime,
				Fields:      tt.fields.Fields,
				ErrorCode:   tt.fields.ErrorCode,
				RequestID:   tt.fields.RequestID,
				IsLogged:    tt.fields.IsLogged,
			}

			_ = l.GetLatency()
		})
	}
}

func TestLogInfoType_SetTime(t *testing.T) {
	type fields struct {
		Name        string
		Request     interface{}
		Response    interface{}
		Err         error
		Latency     string
		requestTime time.Time
		Fields      Fields
		ErrorCode   int
		RequestID   string
		IsLogged    bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "Test Set Time",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LogInfoType{
				Name:        tt.fields.Name,
				Request:     tt.fields.Request,
				Response:    tt.fields.Response,
				Err:         tt.fields.Err,
				Latency:     tt.fields.Latency,
				requestTime: tt.fields.requestTime,
				Fields:      tt.fields.Fields,
				ErrorCode:   tt.fields.ErrorCode,
				RequestID:   tt.fields.RequestID,
				IsLogged:    tt.fields.IsLogged,
			}
			l.SetTime()
		})
	}
}

func TestLogInfoType_String(t *testing.T) {
	type fields struct {
		Name        string
		Request     interface{}
		Response    interface{}
		Err         error
		Latency     string
		requestTime time.Time
		Fields      Fields
		ErrorCode   int
		RequestID   string
		IsLogged    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := LogInfoType{
				Name:        tt.fields.Name,
				Request:     tt.fields.Request,
				Response:    tt.fields.Response,
				Err:         tt.fields.Err,
				Latency:     tt.fields.Latency,
				requestTime: tt.fields.requestTime,
				Fields:      tt.fields.Fields,
				ErrorCode:   tt.fields.ErrorCode,
				RequestID:   tt.fields.RequestID,
				IsLogged:    tt.fields.IsLogged,
			}
			if got := l.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogging(t *testing.T) {
	type args struct {
		ctx      context.Context
		name     string
		isLogged bool
		request  interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  context.Context
		want1 *LogInfoType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Logging(tt.args.ctx, tt.args.name, tt.args.isLogged, tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Logging() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Logging() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInfofCtx(t *testing.T) {
	type args struct {
		ctx    context.Context
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InfofCtx(tt.args.ctx, tt.args.format, tt.args.args...)
		})
	}
}

func Test_printInfo(t *testing.T) {
	type args struct {
		ctx         context.Context
		skipCounter int
		logType     string
		format      string
		args        []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printInfo(tt.args.ctx, tt.args.skipCounter, tt.args.logType, tt.args.format, tt.args.args...)
		})
	}
}

func Test_requestTracer(t *testing.T) {
	type args struct {
		ctx     context.Context
		logInfo *LogInfoType
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := requestTracer(tt.args.ctx, tt.args.logInfo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("requestTracer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slackNotification(t *testing.T) {
	type args struct {
		log      *LogInfoType
		request  interface{}
		response interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slackNotification(tt.args.log, tt.args.request, tt.args.response)
		})
	}
}
