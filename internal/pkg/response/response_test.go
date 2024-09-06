package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-app/internal/pkg/errs"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(middleware []gin.HandlerFunc, routers ...func(r gin.IRoutes) gin.IRoutes) (*gin.Engine, *httptest.ResponseRecorder) {
	r := gin.Default()

	g := r.Group("/test")
	if middleware != nil {
		g.Use(middleware...)
	}
	for _, router := range routers {
		router(g)
	}

	w := httptest.NewRecorder()

	return r, w
}

func TestResponseFull(t *testing.T) {
	Init(Opts{
		WithStackTrace: true,
	})
	tests := []struct {
		name       string
		method     string
		path       string
		handler    func(r gin.IRoutes) gin.IRoutes
		middleware []gin.HandlerFunc
		wantStatus int
	}{
		{
			name:   "Test Data Response",
			method: "GET",
			path:   "/test",
			middleware: []gin.HandlerFunc{
				Middleware,
			},
			handler: func(r gin.IRoutes) gin.IRoutes {
				return r.GET("", func(c *gin.Context) {
					DataResponse(c, gin.H{"success": true})
				})
			},
			wantStatus: http.StatusOK,
		}, {
			name:       "Test Data Response w/o Midleware",
			method:     "GET",
			path:       "/test",
			middleware: []gin.HandlerFunc{},
			handler: func(r gin.IRoutes) gin.IRoutes {
				return r.GET("", func(c *gin.Context) {
					DataResponse(c, gin.H{"success": true})
				})
			},
			wantStatus: http.StatusOK,
		}, {
			name:   "Test Data Response Process Time Tampered",
			method: "GET",
			path:   "/test",
			middleware: []gin.HandlerFunc{
				func(c *gin.Context) {
					c.Set(ProcessingTimeKey, "lolol")
					c.Next()
				},
			},
			handler: func(r gin.IRoutes) gin.IRoutes {
				return r.GET("", func(c *gin.Context) {
					DataResponse(c, gin.H{"success": true})
				})
			},
			wantStatus: http.StatusOK,
		}, {
			name:   "Test Error Resp",
			method: "GET",
			path:   "/test",
			middleware: []gin.HandlerFunc{
				Middleware,
			},
			handler: func(r gin.IRoutes) gin.IRoutes {
				return r.GET("", func(c *gin.Context) {
					Err(c, WrapErr(WrapErrCode(errors.New("some err"), NotFoundCode, "wrap")))
				})
			},
			wantStatus: http.StatusNotFound,
		}, {
			name:   "Test Error Resp Non Wrap",
			method: "GET",
			path:   "/test",
			middleware: []gin.HandlerFunc{
				Middleware,
			},
			handler: func(r gin.IRoutes) gin.IRoutes {
				return r.GET("", func(c *gin.Context) {
					Err(c, errors.New("some err"))
				})
			},
			wantStatus: http.StatusInternalServerError,
		}, {
			name:   "Test New Error Resp",
			method: "GET",
			path:   "/test",
			middleware: []gin.HandlerFunc{
				Middleware,
			},
			handler: func(r gin.IRoutes) gin.IRoutes {
				return r.GET("", func(c *gin.Context) {
					Err(c, WrapErr(NewError("some err", NotFoundCode)))
				})
			},
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w := setupRouter(tt.middleware, tt.handler)
			req, _ := http.NewRequest(tt.method, tt.path, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestGetRequestID(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRequestID(tt.args.c); got != tt.want {
				t.Errorf("GetRequestID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetProcessingTime(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetProcessingTime(tt.args.c); got != tt.want {
				t.Errorf("GetProcessingTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildErrResponse(t *testing.T) {
	type args struct {
		requestID      string
		processingTime string
		err            error
		data           []interface{}
	}
	err1 := errors.New("err")
	tests := []struct {
		name  string
		args  args
		want1 int
		want2 Response
	}{
		{
			name: "success",
			args: args{
				requestID:      "req-id",
				processingTime: "0.01ns",
				err:            wrapErrCode(err1, BadRequestErrCode),
				data:           nil,
			},
			want1: BadRequestErrCode.HTTPCode(),
			want2: Response{
				RequestID:      "req-id",
				Code:           BadRequestErrCode.Code(),
				ProcessingTime: "0.01ns",
				Data:           nil,
				Reason:         BadRequestErrCode.userMsg,
				Error:          BadRequestErrCode.devMsg,
			},
		},
		{
			name: "success - errs v2",
			args: args{
				requestID:      "req-id",
				processingTime: "0.01ns",
				err:            errs.NewErrBadRequest(),
				data:           nil,
			},
			want1: BadRequestErrCode.HTTPCode(),
			want2: Response{
				RequestID:      "req-id",
				Code:           errs.ErrTypeBadRequest.Code,
				ProcessingTime: "0.01ns",
				Data:           nil,
				Reason:         errs.ErrTypeBadRequest.UserMsg(),
				Error:          errs.ErrTypeBadRequest.PublicMsg(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := BuildErrResponse(tt.args.requestID, tt.args.processingTime, tt.args.err, tt.args.data...)
			got2.ErrorDetails = nil

			if got1 != tt.want1 {
				t.Errorf("buildErrResponse() got1 = %v, want %v", got1, tt.want1)
			}
			got2.err = nil
			assert.Equal(t, tt.want2, got2)
		})
	}
}
