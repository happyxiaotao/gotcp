# gotcp
基础的TCP底层框架，暴露三个事件回调接口

# 使用示例

```go
package main

import (
	"log"
	"tttt/gotcp"
)

type Handler struct {
	server gotcp.Server
}

func (s *Handler) OnStart(session gotcp.Session) error {
	log.Printf("OnStart,sessionId:%v", session.SessionId())
	return nil
}

func (s *Handler) OnData(session gotcp.Session, data []byte) error {
	log.Printf("OnData,sessionId:%v,data:%s", session.SessionId(), string(data))

	if session.SessionId() == gotcp.SessionId(2) {
		s.server.Stop()
	}
	return nil
}
func (s *Handler) OnClose(session gotcp.Session) {
	log.Printf("OnClose,sessionId:%v", session.SessionId())
}

func main() {
	var handler Handler
	var address = "127.0.0.1:9991"
	server := gotcp.NewServer(&handler, address)
	handler.server = server
	log.Printf("listen addr:%v", address)
	err := server.Run()
	log.Printf("Run err:%v", err)
}


```
