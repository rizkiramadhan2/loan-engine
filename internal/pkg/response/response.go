package response

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var enableStackTrace bool

// Init init pkg
func EnableStackTrace(withStackTrace bool) {
	enableStackTrace = withStackTrace
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

	code Code
}

// Data write generic data response
func Data(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		RequestID:      GetRequestID(c),
		Data:           data,
		Code:           SuccessCode.code,
		ProcessingTime: GetProcessingTime(c),
	})
}

func buildErr(requestID, processingTime string, err error, data ...interface{}) Response {
	resp := Response{
		RequestID:      requestID,
		ProcessingTime: processingTime,
	}

	e := WrapErr(err)
	resp.ErrorDetails = &ErrorDetails{
		ErrorMsg: e.Error(),
	}
	if e.data != nil {
		resp.Data = data
	}
	if enableStackTrace {
		resp.ErrorDetails.StackTrace = e.StackTrace()
	}

	if len(data) > 0 {
		resp.Data = data
	}

	resp.Reason = e.Readable()
	resp.Code = e.code.code
	resp.code = e.code
	resp.Error = e.code.devMsg

	return resp
}

// ErrCode write generic error response
func ErrCode(c *gin.Context, code Code, data ...interface{}) {
	Err(c, NewError(code.devMsg, code), data...)
}

// Err write generic error response
func Err(c *gin.Context, err error, data ...interface{}) {
	reqID := GetRequestID(c)
	pTime := GetProcessingTime(c)
	resp := buildErr(reqID, pTime, err, data...)
	c.AbortWithStatusJSON(resp.code.HTTPCode(), resp)
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
