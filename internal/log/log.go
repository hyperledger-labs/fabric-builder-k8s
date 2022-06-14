// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"fmt"
	"log"
	"os"
)

type logContextKeyType string

const cmdKey logContextKeyType = "cmd"
const debugKey logContextKeyType = "debug"
const pidKey logContextKeyType = "pid"

type minimalLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type CmdLogger struct {
	infoLogger  minimalLogger
	debugLogger minimalLogger
}

type nilLogger struct{}

func (nl *nilLogger) Print(v ...interface{}) {
	// do nothing
}

func (nl *nilLogger) Printf(format string, v ...interface{}) {
	// do nothing
}

func (nl *nilLogger) Println(v ...interface{}) {
	// do nothing
}

const (
	flags = 0
)

func New(ctx context.Context) *CmdLogger {
	cmd, _ := CmdFromContext(ctx)
	pid, _ := PidFromContext(ctx)

	infoPrefix := fmt.Sprintf("%s [%v]: ", cmd, pid)
	infoLogger := log.New(os.Stderr, infoPrefix, flags)

	var debugLogger minimalLogger
	if DebugFromContext(ctx) {
		debugPrefix := fmt.Sprintf("%s [%v] DEBUG: ", cmd, pid)
		debugLogger = log.New(os.Stderr, debugPrefix, flags)
	} else {
		debugLogger = &nilLogger{}
	}

	cl := &CmdLogger{
		infoLogger:  infoLogger,
		debugLogger: debugLogger,
	}

	return cl
}

// NewCmdContext returns a new Context with program name, process id, and
// debug values
func NewCmdContext(ctx context.Context, debug bool) context.Context {
	cmdValue := os.Args[0]
	cmdContext := context.WithValue(ctx, cmdKey, cmdValue)

	pidValue := os.Getpid()
	cmdContext = context.WithValue(cmdContext, pidKey, pidValue)

	cmdContext = context.WithValue(cmdContext, debugKey, debug)

	return cmdContext
}

// CmdFromContext returns the program name value from the provided Context
func CmdFromContext(ctx context.Context) (string, bool) {
	cmd, ok := ctx.Value(cmdKey).(string)
	return cmd, ok
}

// PidFromContext returns the process ID value from the provided Context
func PidFromContext(ctx context.Context) (int, bool) {
	pid, ok := ctx.Value(pidKey).(int)
	return pid, ok
}

// DebugFromContext returns if debug is enabled in the provided Context
func DebugFromContext(ctx context.Context) bool {
	if d, ok := ctx.Value(debugKey).(bool); ok {
		return d
	}
	return false
}

func (cl *CmdLogger) Print(v ...interface{}) {
	cl.infoLogger.Print(v...)
}

func (cl *CmdLogger) Printf(format string, v ...interface{}) {
	cl.infoLogger.Printf(format, v...)
}

func (cl *CmdLogger) Println(v ...interface{}) {
	cl.infoLogger.Println(v...)
}

func (cl *CmdLogger) Debug(v ...interface{}) {
	cl.debugLogger.Print(v...)
}

func (cl *CmdLogger) Debugf(format string, v ...interface{}) {
	cl.debugLogger.Printf(format, v...)
}

func (cl *CmdLogger) Debugln(v ...interface{}) {
	cl.debugLogger.Println(v...)
}
