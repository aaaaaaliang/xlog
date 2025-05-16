# cil-logger

一个轻量级、高可读性、无依赖的 Go 日志库，支持彩色输出、结构化日志（JSON）、链式调用、字段排序、trace_id 注入、线程安全等功能。

---

## ✨ 功能特性

- ✅ 支持日志级别（DEBUG、INFO、ERROR）
- ✅ 支持输出格式（Text / JSON）
- ✅ 自动识别生产环境（ENV=prod）切换为 JSON 输出
- ✅ 彩色控制台输出（Text 模式下）
- ✅ 支持字段增强（WithField / WithFields）
- ✅ 字段排序输出（保证日志结构一致）
- ✅ 支持文件输出 / 多路输出（控制台 + 文件）
- ✅ 线程安全，支持高并发
- ✅ 可手动注入 trace_id（支持 base36 编码）
- ✅ 自动记录代码调用位置（文件名 + 行号）

---

## 📦 安装使用

将 `logger/` 文件夹拷贝到你的项目中作为内部模块使用。

---

## 🧪 快速使用示例

```go
package main

import (
    "time"
)

func main() {
    traceID := logger.NewTraceID()

    log := logger.New("OrderService").
        WithField("trace_id", traceID)

    start := time.Now()
    time.Sleep(120 * time.Millisecond)
    log.WithField("duration", time.Since(start)).
        Info("订单处理完成")
}

```

---

## 🔧 配置日志格式

```go
log := logger.New("UserService")
log.SetFormat(logger.FormatJSON) // 设置为 JSON 格式输出
```

也可以通过环境变量自动切换：

```bash
ENV=prod go run main.go  # 自动使用 JSON 格式
```

---

## 🗃️ 输出到文件

```go
file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err != nil {
    panic(err)
}
log.SetOutput(io.MultiWriter(os.Stdout, file)) // 同时输出到终端和文件
```

---

## 📍 获取 trace_id（默认 base36 编码）

```go
traceID := logger.NewTraceID() // 如：kxz0w1b4i80
```

---

## 🧠 注意事项

- `logger.New(...)` 返回的是接口 `Logger`，支持链式调用。
- JSON 格式输出时字段统一放入 `fields` 字段中，方便日志采集。
- 字段输出自动排序，提升一致性和可读性。

---

## 🧩 适合场景

- 中小型项目快速接入结构化日志
- 替代标准库 `log` 的升级版本
- 配合 logstash / filebeat 收集日志
- 多环境灵活切换（Text / JSON）

---

## 📄 License

MIT
