package main

import (
	"fmt"
	"strings"
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
	SetLevel(level Level)
}

type cilLogger struct {
	name   string
	fields map[string]interface{}
	level  Level
}

func New(name string) Logger {
	return &cilLogger{
		name:   name,
		fields: make(map[string]interface{}),
		level:  LevelDebug,
	}
}

func (l *cilLogger) clone() *cilLogger {
	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	return &cilLogger{
		name:   l.name,
		fields: newFields,
		level:  l.level,
	}
}

func (l *cilLogger) WithFields(fields map[string]interface{}) Logger {
	clone := l.clone()
	for k, v := range fields {
		clone.fields[k] = v
	}
	return clone
}

func (l *cilLogger) SetLevel(level Level) {
	l.level = level
}

func (l *cilLogger) shouldLog(level Level) bool {
	return level >= l.level
}

func (l *cilLogger) formatLog(level Level, msg string) string {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	color := level.Color()
	reset := "\033[0m"

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s[%s] [%s] [%s]", color, level.String(), timeStr, l.name))
	for k, v := range l.fields {
		sb.WriteString(fmt.Sprintf(" %s=%v", k, v))
	}
	sb.WriteString(fmt.Sprintf(" | %s%s\n", msg, reset))
	return sb.String()
}

func (l *cilLogger) log(level Level, msg string) {
	if !l.shouldLog(level) {
		return
	}
	fmt.Print(l.formatLog(level, msg))
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
	log := New("OrderService")
	log.SetLevel(LevelDebug)

	log.Debug("调试信息")
	log.Info("服务启动成功")

	userLog := log.WithFields(map[string]interface{}{
		"user_id": 42,
		"ip":      "192.168.1.1",
	})

	userLog.Info("用户登录成功")
	userLog.Error("数据库连接失败")

	fmt.Println(LevelInfo)         // 会打印 INFO
	fmt.Printf("%s\n", LevelError) // 会打印 ERROR
	fmt.Printf("%d\n", LevelError) // 会打印 2（LevelError 的底层 int 值）

}
