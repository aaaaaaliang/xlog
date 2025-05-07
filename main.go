package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func (l Level) Color() string {
	switch l {
	case LevelDebug:
		return "\033[36m" // 青色
	case LevelInfo:
		return "\033[32m" // 绿色
	case LevelError:
		return "\033[31m" // 红色
	default:
		return "\033[0m"
	}
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Error(msg string)
	WithFields(fields map[string]interface{}) Logger
	WithField(key string, value interface{}) Logger
	SetLevel(level Level)
	SetOutput(w io.Writer)
}

type cilLogger struct {
	name   string
	fields map[string]interface{}
	level  Level
	output io.Writer
	lock   sync.RWMutex
}

func New(name string) Logger {
	return &cilLogger{
		name:   name,
		fields: make(map[string]interface{}),
		level:  LevelDebug,
		output: os.Stdout,
	}
}

func (l *cilLogger) clone() *cilLogger {
	l.lock.RLock()
	defer l.lock.RUnlock()

	newFields := make(map[string]interface{}, len(l.fields))
	for k, v := range l.fields {
		newFields[k] = v
	}

	return &cilLogger{
		name:   l.name,
		fields: newFields,
		level:  l.level,
		output: l.output,
	}
}

func (l *cilLogger) WithFields(fields map[string]interface{}) Logger {
	clone := l.clone()
	for k, v := range fields {
		clone.fields[k] = v
	}
	return clone
}

func (l *cilLogger) WithField(key string, value interface{}) Logger {
	return l.WithFields(map[string]interface{}{key: value})
}

func (l *cilLogger) SetLevel(level Level) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.level = level
}

func (l *cilLogger) shouldLog(level Level) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return level >= l.level
}

func (l *cilLogger) SetOutput(w io.Writer) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.output = w
}

func (l *cilLogger) formatLog(level Level, msg string) string {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	color := level.Color()
	reset := "\033[0m"

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s[%s] [%s] [%s]", color, level.String(), timeStr, l.name))

	// 🚀 对 key 排序
	keys := make([]string, 0, len(l.fields))
	for k := range l.fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		sb.WriteString(fmt.Sprintf(" %s=%v", k, l.fields[k]))
	}

	sb.WriteString(fmt.Sprintf(" | %s%s\n", msg, reset))
	return sb.String()
}

func (l *cilLogger) log(level Level, msg string) {
	if !l.shouldLog(level) {
		return
	}

	// 获取输出锁
	l.lock.RLock()
	output := l.output
	l.lock.RUnlock()

	_, _ = output.Write([]byte(l.formatLog(level, msg)))
}

func (l *cilLogger) Debug(msg string) {
	l.log(LevelDebug, msg)
}

func (l *cilLogger) Info(msg string) {
	l.log(LevelInfo, msg)
}

func (l *cilLogger) Error(msg string) {
	l.log(LevelError, msg)
}

func main() {
	log := New("UserService")

	for i := 0; i < 100; i++ {
		go func(i int) {
			log.WithField("user", i).Info("处理完毕")
		}(i)
	}

	time.Sleep(time.Second)
}
