package gotcp

import "sync/atomic"

type SessionId uint64

var globalSessionId uint64

func generateSessionId() SessionId {
	return SessionId(atomic.AddUint64(&globalSessionId, uint64(1)))
}
