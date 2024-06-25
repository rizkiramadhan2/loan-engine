package response

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// ProcessingTimeKey start time key
	ProcessingTimeKey = "trace-request-start_time"
	// RequestIDKey request id key
	RequestIDKey = "trace-request_id"
)

// Middleware start processing time and request id inside gin.Context
func Middleware(c *gin.Context) {
	c.Set(ProcessingTimeKey, time.Now())
	c.Set(RequestIDKey, uuid.New().String())

	c.Next()
}
