package main

import (
	"cil-log/logger"
	"time"
)

func main() {
	traceID := logger.NewTraceID()

	log := logger.New("OrderService").WithField("trace_id", traceID)

	start := time.Now()
	time.Sleep(150 * time.Millisecond)
	log.WithField("duration", time.Since(start)).Info("订单处理完成")
}
