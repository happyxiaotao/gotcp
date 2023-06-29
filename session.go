package gotcp

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Session interface {
	Start()                    // 会话开启
	Close(err ...any)          // 关闭会话
	Write([]byte) (int, error) // 写入数据
	SessionId() SessionId      // 返回SessionId
}

type session struct {
	sessionId SessionId
	conn      net.Conn           // 底层网络连接
	handler   EventHandler       // 事件回调处理
	isClosed  atomic.Bool        // 确定连接是否关闭，避免连接关闭后，还是发送数据
	ctx       context.Context    // 当前会话的Context，用来处理此会话关联或派生的所有协程。
	cancel    context.CancelFunc // 当前会话的Context对应的cancel函数，与ctx配套使用
}

func NewSession(conn net.Conn, handler EventHandler) Session {
	ctx, cancel := context.WithCancel(context.Background())
	return &session{
		sessionId: generateSessionId(),
		conn:      conn,
		handler:   handler,
		isClosed:  atomic.Bool{},
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (s *session) Start() {
	if s.handler != nil {
		err := s.handler.OnStart(s)
		if err != nil {
			s.Close(err)
			return // 失败就直接退出函数
		}
	}
	s.readloop(s.ctx)
}

func (s *session) Close(err ...any) {
	if s.isClosed.Load() {
		return
	}
	log.Printf("Close err:%v", err)
	s.isClosed.Store(true)
	s.cancel()
	s.conn.Close()
	if s.handler != nil {
		s.handler.OnClose(s)
	}
}

// readloop协程
func (s *session) readloop(ctx context.Context) {
	var reader *bufio.Reader = bufio.NewReader(s.conn)
	var buf = make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			// err := ctx.Err()
			// log.Printf("session:readloop ctx.Done,sessionId:%v,err:%v", s.SessionId(), err)
			return
		default:
			// 指定从buffer中读取数据的最大容量
			// 从buffer中读取数据并保存到buf中，n代表实际返回的数据大小
			n, err := reader.Read(buf)
			if err != nil {
				s.Close(err)
				break
			}
			if n == 0 {
				continue
			}
			if s.handler != nil {
				err := s.handler.OnData(s, buf[:n])
				if err != nil {
					s.Close(err)
				}
			}
		}
	}
}

func (s *session) Write(data []byte) (int, error) {
	if !s.isClosed.Load() {
		return s.conn.Write(data)
	} else {
		return 0, fmt.Errorf("连接已关闭")
	}
}

func (s *session) SessionId() SessionId { return s.sessionId }
