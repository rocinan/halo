package halo

import (
	"io"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type WebSocket struct {
	Fd int
	Ln net.Listener
}

func (t *WebSocket) Listen(sa *Addr) (SConn, error) {
	la, err := net.ListenTCP("tcp", &net.TCPAddr{IP: sa.IP[:], Port: sa.Port})
	if err != nil {
		return nil, err
	}
	return &WebSocket{
		Ln: la,
	}, nil
}

func (t *WebSocket) Accept() (Conn, error) {
	conn, err := t.Ln.Accept()
	if err != nil {
		return nil, err
	}
	cn, _ := conn.(*net.TCPConn)
	return &WSConn{
		Conn:   cn,
		Reader: wsutil.NewReader(conn, ws.StateServerSide),
		Writer: wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText),
	}, nil
}

func (t *WebSocket) Dial(sa *Addr) (Conn, error) {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: sa.IP[:], Port: sa.Port})
	if err != nil {
		return nil, err
	}
	return &WSConn{
		Conn:   conn,
		Reader: wsutil.NewReader(conn, ws.StateServerSide),
		Writer: wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText),
	}, nil
}

func (t *WebSocket) Close() error {
	return t.Ln.Close()
}

type WSConn struct {
	Fd     int
	Conn   *net.TCPConn
	Reader *wsutil.Reader
	Writer *wsutil.Writer
}

func (c *WSConn) Read(b []byte) (n int, err error) {
	return io.ReadFull(c.Reader, b)
}

func (c *WSConn) Write(b []byte) (n int, err error) {
	return c.Writer.Write(b)
}

func (c *WSConn) Close() error {
	return c.Conn.Close()
}
