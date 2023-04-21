package zap

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// Decider function defines rules for suppressing any interceptor logs
type Decider func(requestUri string) bool

// DefaultDeciderMethod is the default implementation of decider to see if you should log the call
// by default this if always true so all calls are logged
func DefaultDeciderMethod(requestUri string) bool {
	return true
}

// DurationToField function defines how to produce duration fields for logging
type DurationToField func(duration time.Duration) zapcore.Field
