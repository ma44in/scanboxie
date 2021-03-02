package device

// Thx to https://github.com/gvalkov/golang-evdev
//
// Just copied needed parts from gvalkov/golang-evdev to
// avoid cgo usage.

import (
	"syscall"
	"unsafe"
)

// InputEvent defines event struct
type InputEvent struct {
	Time  syscall.Timeval // time in seconds since epoch at which event occurred
	Type  uint16          // event type - one of ecodes.EV_*
	Code  uint16          // event code related to the event type
	Value int32           // event value related to the event type
}

var eventsize = int(unsafe.Sizeof(InputEvent{}))
