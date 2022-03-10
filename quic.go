package halo

import (
	"context"
	"crypto/tls"
	"net"

	quic "github.com/lucas-clemente/quic-go"
)

type Quic struct {
	Fd int
	Ln quic.Listener
}

func (qc *Quic) Listen(sa *Addr) (SConn, error) {
	cn, err := net.ListenUDP("udp", &net.UDPAddr{IP: sa.IP[:], Port: sa.Port})
	if err != nil {
		return nil, err
	}
	la, err := quic.Listen(cn, GenerateTLSConfig(), nil)
	if err != nil {
		return nil, err
	}
	return &Quic{
		Ln: la,
	}, nil
}

func (qc *Quic) Accept() (Conn, error) {
	sess, err := qc.Ln.Accept(context.Background())
	if err != nil {
		return nil, err
	}
	stream, err := sess.AcceptStream(context.Background())
	if err != nil {
		return nil, err
	}
	return &QuicConn{
		Fd:   qc.Fd,
		Conn: stream,
	}, nil
}

func (qc *Quic) Dial(sa *Addr) (Conn, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"secret"},
	}
	la, err := net.ListenUDP("udp", nil)
	if err != nil {
		return nil, err
	}
	sess, err := quic.Dial(la, &net.UDPAddr{IP: sa.IP[:], Port: sa.Port}, "", tlsConf, nil)
	if err != nil {
		return nil, err
	}
	stream, err := sess.OpenStreamSync(context.Background())
	if err != nil {
		return nil, err
	}
	return &QuicConn{
		La:   la,
		Conn: stream,
	}, nil
}

func (qc *Quic) Close() error {
	return qc.Ln.Close()
}

type QuicConn struct {
	Fd   int
	La   net.PacketConn
	Conn quic.Stream
}

func (c *QuicConn) Read(b []byte) (n int, err error) {
	return c.Conn.Read(b)
}

func (c *QuicConn) Write(b []byte) (n int, err error) {
	return c.Conn.Write(b)
}

func (c *QuicConn) Close() error {
	return c.Conn.Close()
}
