package response

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"simple-app/internal/pkg/errs"
	"simple-app/internal/pkg/log"

	"github.com/gin-gonic/gin"
)

var enableStackTrace bool

type Opts struct {
	WithStackTrace  bool
	GenericErrorMsg string
}

// Init init pkg
func Init(opts Opts) {
	enableStackTrace = opts.WithStackTrace
	if opts.GenericErrorMsg != "" {
		genericErrorMsg = opts.GenericErrorMsg
	}
	setDefaults()
}

// ErrorDetails for staging and developmnent only
type ErrorDetails struct {
	ErrorMsg   string   `json:"error_msg,omitempty"`
	StackTrace []string `json:"stack_trace,omitempty"`
}

// Response basic response
type Response struct {
	RequestID      string      `json:"request_id"`
	Code           string      `json:"code"`
	ProcessingTime string      `json:"processing_time,omitempty"`
	Data           interface{} `json:"data,omitempty"`

	Reason string `json:"reason,omitempty"`

	Error        string        `json:"error,omitempty"`
	ErrorDetails *ErrorDetails `json:"error_details,omitempty"`

	err error
}

func (r *Response) LogCtx(ctx context.Context) {
	if r.Error != "" {
		det, _ := r.GetErrorDetails()
		log.ErrorfCtx(ctx, "reqID: %s, code: %s, err: %s, trace:\n%s", r.RequestID, r.Code, det.ErrorMsg, strings.Join(det.StackTrace, "\n"))
	} else {
		log.InfofCtx(ctx, "reqID: %s, code: %s, data: %+v", r.RequestID, r.Code, r.Data)
	}
}
func (r *Response) Log() {
	if r.Error != "" {
		det, _ := r.GetErrorDetails()
		log.Errorf("reqID: %s, code: %s, err: %s, trace:\n%s", r.RequestID, r.Code, det.ErrorMsg, strings.Join(det.StackTrace, "\n"))
	} else {
		log.Infof("reqID: %s, code: %s, data: %+v", r.RequestID, r.Code, r.Data)
	}
}

func (r *Response) GetErrorDetails() (ErrorDetails, any) {
	e := WrapErr(r.err)
	det := ErrorDetails{
		ErrorMsg:   e.Error(),
		StackTrace: e.StackTrace(),
	}
	return det, e.data
}

// DataResponse write generic data response
func DataResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		RequestID:      GetRequestID(c),
		Data:           data,
		Code:           SuccessCode.code,
		ProcessingTime: GetProcessingTime(c),
	})
}

func BuildErrResponse(requestID, processingTime string, err error, data ...interface{}) (int, Response) {
	var er *errs.Err
	if errors.As(err, &er) {
		return buildErrResponseV2(err, requestID, processingTime, data...)
	}

	e := WrapErr(err)

	resp := Response{
		RequestID:      requestID,
		Code:           e.code.code,
		ProcessingTime: processingTime,

		Reason: e.Readable(),
		Error:  e.code.devMsg,
	}

	resp.err = e
	if enableStackTrace {
		det, data := resp.GetErrorDetails()
		resp.ErrorDetails = &det
		if data != nil {
			resp.Data = e.data
		}
	}

	if len(data) > 0 {
		resp.Data = data
	}

	return e.code.HTTPCode(), resp
}

func buildErrResponseV2(e error, requestID, processingTime string, data ...interface{}) (int, Response) {
	var err *errs.Err
	errors.As(e, &err)

	var httpErr *errs.HTTPStatus
	if !errors.As(e, &httpErr) {
		errors.As(errs.ErrTypeInternalErr.Copy(), &httpErr)
	}

	resp := Response{
		RequestID:      requestID,
		Code:           httpErr.Code,
		ProcessingTime: processingTime,

		Reason: err.UserMsg(),
		Error:  err.PublicMsg(),
	}

	resp.err = err
	if enableStackTrace {
		resp.ErrorDetails = &ErrorDetails{
			ErrorMsg:   err.Error(),
			StackTrace: err.StackTrace(),
		}
	}

	if len(data) > 0 {
		resp.Data = data
	}

	return httpErr.HTTPStatusCode(), resp
}

// ErrCode write generic error response
func ErrCode(c *gin.Context, code Code, data ...interface{}) {
	Err(c, NewError(code.devMsg, code), data...)
}

// Err write generic error response
func Err(c *gin.Context, err error, data ...interface{}) *Response {
	reqID := GetRequestID(c)
	pTime := GetProcessingTime(c)
	statusCode, resp := BuildErrResponse(reqID, pTime, err, data...)
	c.AbortWithStatusJSON(statusCode, resp)
	return &resp
}

// GetRequestID get request id from context
func GetRequestID(c *gin.Context) string {
	v, ok := c.Get(RequestIDKey)
	if !ok {
		return ""
	}

	return fmt.Sprintf("%v", v)
}

// GetProcessingTime get request processing time
func GetProcessingTime(c *gin.Context) string {
	v, ok := c.Get(ProcessingTimeKey)
	if !ok {
		return ""
	}

	startTime, ok := v.(time.Time)
	if !ok {
		return ""
	}

	// take the float value from Seconds() instead returning from Milliseconds()
	return fmt.Sprintf("%.2fms", time.Since(startTime).Seconds()*1000)
}
