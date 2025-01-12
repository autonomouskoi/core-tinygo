package core

import (
	"errors"
	"fmt"
)

// Log a message on the host. If args are present there should be an even
// number of them. The even numbered args must be string keys. The odd
// numbered args must be an integer type, float type, string, or bool.
// It's slightly more convenient to use LogError(), LogInfo, etc
func Log(level LogLevel, message string, args ...any) error {
	if len(args)%2 != 0 {
		return errors.New("arg count must be even")
	}
	logArgs := make([]*LogSendRequest_Arg, len(args)/2)
	for i := 1; i < len(args); i += 2 {
		key, ok := args[i-1].(string)
		if !ok {
			return fmt.Errorf("arg %d not a key string", i)
		}
		arg := &LogSendRequest_Arg{Key: key}
		switch v := args[i].(type) {
		case float32:
			arg.Value = &LogSendRequest_Arg_Double{Double: float64(v)}
		case float64:
			arg.Value = &LogSendRequest_Arg_Double{Double: v}
		case int8:
			arg.Value = &LogSendRequest_Arg_Int64{Int64: int64(v)}
		case int16:
			arg.Value = &LogSendRequest_Arg_Int64{Int64: int64(v)}
		case int32:
			arg.Value = &LogSendRequest_Arg_Int64{Int64: int64(v)}
		case int64:
			arg.Value = &LogSendRequest_Arg_Int64{Int64: v}
		case uint:
			arg.Value = &LogSendRequest_Arg_Int64{Int64: int64(v)}
		case uint8:
			arg.Value = &LogSendRequest_Arg_Int64{Int64: int64(v)}
		case uint16:
			arg.Value = &LogSendRequest_Arg_Int64{Int64: int64(v)}
		case uint32:
			arg.Value = &LogSendRequest_Arg_Int64{Int64: int64(v)}
		case uint64:
			arg.Value = &LogSendRequest_Arg_Int64{Int64: int64(v)}
		case string:
			arg.Value = &LogSendRequest_Arg_String_{String_: v}
		case bool:
			arg.Value = &LogSendRequest_Arg_Bool{Bool: v}
		default:
			return fmt.Errorf("unhandled type %T: %v", v, v)
		}
		logArgs[i/2] = arg
	}

	msg := &BusMessage{
		Type: int32(ExternalMessageType_LOG_SEND_REQ),
	}
	req := &LogSendRequest{
		Level:   level,
		Message: message,
		Args:    logArgs,
	}
	b, err := req.MarshalVT()
	if err != nil {
		return fmt.Errorf("marshalling: %w", err)
	}
	msg.Message = b
	return Send(msg)
}

// LogError logs a message at level ERROR on the host. See the docs for Log()
// for a description of args
func LogError(message string, args ...any) error {
	return Log(LogLevel_ERROR, message, args...)
}

// LogWarn logs a message at level WARN on the host. See the docs for Log()
// for a description of args
func LogWarn(message string, args ...any) error {
	return Log(LogLevel_WARN, message, args...)
}

// LogInfo logs a message at level INFO on the host. See the docs for Log()
// for a description of args
func LogInfo(message string, args ...any) error {
	return Log(LogLevel_INFO, message, args...)
}

// LogDebug logs a message at level DEBUG on the host. See the docs for Log()
// for a description of args
func LogDebug(message string, args ...any) error {
	return Log(LogLevel_DEBUG, message, args...)
}
