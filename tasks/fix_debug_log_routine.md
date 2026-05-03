# Task: Fix Debug Log Routine

## Description
The `debug/debug.go` file implements a background `logRoutine` that listens on `msgChan`, but the public logging functions (`Debug`, `InternalDebug`, `Warn`, etc.) call `write` and `writef` directly instead of sending messages to the channel.

```go
func writeToChan(level LogLevel, msg string) {
	write(level, msg) // Should be: msgChan <- message{level: level, format: msg}
}
```

This makes the `logRoutine` dead code and defeats the purpose of asynchronous logging.

## Goals
- [ ] Update `writeToChan` and `writeToChanF` to actually send messages to `msgChan`.
- [ ] Ensure `Shutdown()` correctly waits for the channel to be drained if necessary, or at least closes it gracefully.
- [ ] Verify that logging still works and is now actually handled by the background goroutine.

## Verification Criteria
- [ ] `logRoutine` is no longer dead code.
- [ ] Logs appear on the console as expected.
