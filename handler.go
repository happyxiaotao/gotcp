package gotcp

type EventHandler interface {
	OnStart(Session) error        // 会话启动回调
	OnData(Session, []byte) error // 会话读取到数据回调
	OnClose(Session)              // 会话关闭回调
}
