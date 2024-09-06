package log

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	goerrors "github.com/go-errors/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type ctxKey string

const (
	RequestIDKey   string = "trace-request_id" // TODO: gin context is annoying. can't use custom type
	RequestNameKey ctxKey = "trace-request_name"
)

// Fields logrus.Fields
type Fields logrus.Fields

type LogInfoType struct {
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

// String print string
func (l LogInfoType) String() string {
	return fmt.Sprintf("Err=%s, ErrorCode=%d, Name=%s, Latency=%s, RequestId=%s, requestTime=%s, Request=%+v, Fields=%+v, Response=%+v",
		l.Err, l.ErrorCode, l.Name, l.Latency, l.RequestID, l.requestTime.Format("2006-01-02 15:04:05"), l.Request, l.Fields, l.Response)
}

// SetTime set time
func (l *LogInfoType) SetTime() {
	l.requestTime = time.Now()
}

// GetLatency get latency
func (l *LogInfoType) GetLatency() string {
	latency := time.Since(l.requestTime).Seconds() * 1000
	l.Latency = fmt.Sprintf("%.2f ms", latency)
	return l.Latency
}

// HandleDefer handle log defer
func HandleDefer(ctx context.Context, logInfo *LogInfoType) {
	reqStr, _ := json.Marshal(logInfo.Request)
	respStr, _ := json.Marshal(logInfo.Response)

	defer func() {
		var f Fields
		if logInfo.Fields != nil {
			f = logInfo.Fields
		} else {
			f = Fields{
				"method":  logInfo.Name,
				"request": string(reqStr),
				"reply":   string(respStr),
				"err":     logInfo.Err,
			}
		}
		f["latency"] = logInfo.GetLatency()

		if logInfo.IsLogged && logInfo.Err != nil && logInfo.ErrorCode != 401 {
			trace, ok := logInfo.Err.(*goerrors.Error)
			if ok {
				printInfo(ctx, 3, "error", "Error: %+v. Stacktrace:\n%s", logInfo, trace.Stack())
				slackNotification(logInfo, string(reqStr), string(respStr))
			} else {
				printInfo(ctx, 3, "error", "Error: %+v", logInfo)
				slackNotification(logInfo, string(reqStr), string(respStr))
			}
		}
	}()
}

func requestTracer(ctx context.Context, logInfo *LogInfoType) context.Context {
	logInfo.SetTime()
	requestID := uuid.NewV4().String()
	val, ok := ctx.Value(RequestIDKey).(string)
	if ok {
		requestID = val
	}

	logInfo.RequestID = requestID
	ctx = context.WithValue(ctx, RequestIDKey, requestID)
	ctx = context.WithValue(ctx, RequestNameKey, logInfo.Name)

	return ctx
}

func slackNotification(log *LogInfoType, request interface{}, response interface{}) {
	// if slackClient == nil {
	// 	Info("cannot sent slack notification. client not set")
	// 	return
	// }

	// pc, _, _, _ := runtime.Caller(3)

	if log.Err == nil || log.Err == context.Canceled {
		return
	}

	if log.ErrorCode >= 400 && log.ErrorCode != 401 {
		return
	}

	postData := map[string]interface{}{
		"RequestID": log.RequestID,
		"Request":   request,
		"Response":  response,
	}
	if log.Err == context.DeadlineExceeded {
		postData["latency"] = log.GetLatency()
	}

	// slackClient.PostToSlack(nil, log.Err, slack.MessageTypeErrorHandler, string(debug.Stack()), runtime.FuncForPC(pc).Name(), postData)
}

// Logging custom logging
func Logging(ctx context.Context, name string, isLogged bool, request interface{}) (context.Context, *LogInfoType) {
	info := &LogInfoType{Name: name, Request: request, IsLogged: isLogged}

	ctxWithValue := requestTracer(ctx, info)
	return ctxWithValue, info
}

func printInfo(ctx context.Context, skipCounter int, logType string, format string, args ...interface{}) {
	name, _ := ctx.Value(RequestNameKey).(string)
	traceID, _ := ctx.Value(RequestIDKey).(string)

	function := ""
	pc, file, line, ok := runtime.Caller(skipCounter)
	if !ok {
		file = "<???>"
		line = 0
	} else {
		slash := strings.LastIndex(file, "/")
		slash = strings.LastIndex(file[:slash], "/")
		slash = strings.LastIndex(file[:slash], "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}

	params := GetParamsCtx(ctx)
	paramsMsg := ""
	for k, v := range params {
		if paramsMsg != "" {
			paramsMsg += ", "
		}
		paramsMsg += fmt.Sprintf("%v: %+v", k, v)
	}

	// TODO: use a shorter id format for request id
	format = fmt.Sprintf("[%s] [reqID:%s] [%s] %s [%s:%d %s] %s", logType, traceID, name, format, file, line, function, paramsMsg)
	if logType == "error" {
		errLogger.Errorf(format, args...)
		return
	}

	infoLogger.Infof(format, args...)
}

func sprintInfo(ctx context.Context, skipCounter int, logType, format string) string {
	name, _ := ctx.Value(RequestNameKey).(string)
	traceID, _ := ctx.Value(RequestIDKey).(string)

	function := ""
	pc, file, line, ok := runtime.Caller(skipCounter)
	if !ok {
		file = "<???>"
		line = 0
	} else {
		slash := strings.LastIndex(file, "/")
		slash = strings.LastIndex(file[:slash], "/")
		slash = strings.LastIndex(file[:slash], "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()
		function = function[strings.LastIndex(function, "/")+1:]
	}

	params := GetParamsCtx(ctx)
	paramsMsg := ""
	for k, v := range params {
		if paramsMsg != "" {
			paramsMsg += ", "
		}
		paramsMsg += fmt.Sprintf("%v: %+v", k, v)
	}

	// TODO: use a shorter id format for request id
	return fmt.Sprintf("[%s] [reqID:%s] [%s] %s [%s:%d %s] %s", logType, traceID, name, format, file, line, function, paramsMsg)
}

// ErrorfCtx print error ctx
func ErrorfCtx(ctx context.Context, format string, args ...interface{}) {
	format = sprintInfo(ctx, 2, "error", format)
	errLogger.Errorf(format, args...)
}

// InfofCtx print ctx
func InfofCtx(ctx context.Context, format string, args ...interface{}) {
	format = sprintInfo(ctx, 2, "info", format)
	infoLogger.Infof(format, args...)
}

// ErrorfCtxWithSkip print error ctx with trace level skip
func ErrorfCtxWithSkip(ctx context.Context, skip int, format string, args ...interface{}) {
	format = sprintInfo(ctx, 2+skip, "error", format)
	errLogger.Errorf(format, args...)
}

// InfofCtxWithSkip print ctx with trace level skip
func InfofCtxWithSkip(ctx context.Context, skip int, format string, args ...interface{}) {
	format = sprintInfo(ctx, 2+skip, "info", format)
	infoLogger.Infof(format, args...)
}

type logParamsKeyType string

var logParamsKey logParamsKeyType = "log_params"

// AddParamCtx add log parameter to context
func AddParamCtx(ctx context.Context, key string, value interface{}) context.Context {
	params, _ := ctx.Value(logParamsKey).(map[string]interface{})
	newParams := map[string]interface{}{}
	for k, v := range params {
		newParams[k] = v
	}
	newParams[key] = value
	return context.WithValue(ctx, logParamsKey, newParams)
}

// AddParamsCtx add multiple log parameter to context
func AddParamsCtx(ctx context.Context, p map[string]interface{}) context.Context {
	params, _ := ctx.Value(logParamsKey).(map[string]interface{})
	newParams := map[string]interface{}{}
	for k, v := range params {
		newParams[k] = v
	}
	for k, v := range p {
		newParams[k] = v
	}
	return context.WithValue(ctx, logParamsKey, newParams)
}

// SetParamsCtx set log parameter in context
func SetParamsCtx(ctx context.Context, params map[string]interface{}) context.Context {
	return context.WithValue(ctx, logParamsKey, params)
}

// GetParamsCtx get log parameter from context
func GetParamsCtx(ctx context.Context) map[string]interface{} {
	params, _ := ctx.Value(logParamsKey).(map[string]interface{})
	if params == nil {
		params = map[string]interface{}{}
	}
	return params
}
