# gotcp
基础的TCP底层框架，暴露三个事件回调接口

# 使用示例

```go
package main

import (
	"log"
	"tttt/gotcp"
)

type Server struct {
	addr string
}

func (s *Server) OnStart(session gotcp.Session) error {
	log.Printf("OnStart,sessionId:%v", session.SessionId())
	return nil
}

func (s *Server) OnData(session gotcp.Session, data []byte) error {
	log.Printf("OnData,sessionId:%v,data:%s", session.SessionId(), string(data))
	return nil
}
func (s *Server) OnClose(session gotcp.Session) {
	log.Printf("OnClose,sessionId:%v", session.SessionId())
}

func main() {
	server := Server{
		addr: "127.0.0.1:9991",
	}
	log.Printf("listen addr:%v", server.addr)
	err := gotcp.Run(&server, server.addr)
	log.Printf("Run err:%v", err)
}

```
