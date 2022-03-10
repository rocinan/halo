package halo

import (
	"errors"
)

type SConn interface {
	Accept() (Conn, error)
	Close() error
}

type Conn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
}

type PConn interface {
	SendMsg(b []byte) error
	SendMsgTo(b []byte, sa *Addr) error
	RecvMsg(b []byte) (n int, sa *Addr, err error)
	Close() error
}

func Listen(proto, addr string) (SConn, error) {
	sa, err := ResolveAddr(addr)
	if err != nil {
		return nil, errors.New("resolve addr err: " + err.Error())
	}
	switch proto {
	case "ws":
		return new(WebSocket).Listen(sa)
	case "tcp":
		return new(TCP).Listen(sa)
	case "quic":
		return new(Quic).Listen(sa)
	default:
		return nil, errors.New("undefined proto " + proto)
	}
}

func ListenMsg(proto, addr string) (PConn, error) {
	sa, err := ResolveAddr(addr)
	if err != nil {
		return nil, errors.New("resolve addr err: " + err.Error())
	}
	switch proto {
	case "udp":
		return new(UDP).Listen(sa)
	default:
		return nil, errors.New("undefined proto " + proto)
	}
}

func Dial(proto, addr string) (Conn, error) {
	sa, err := ResolveAddr(addr)
	if err != nil {
		return nil, errors.New("resolve addr err: " + err.Error())
	}
	switch proto {
	case "ws":
		return new(WebSocket).Dial(sa)
	case "tcp":
		return new(TCP).Dial(sa)
	case "quic":
		return new(Quic).Dial(sa)
	default:
		return nil, errors.New("undefined proto " + proto)
	}
}

func DialMsg(proto, addr string) (PConn, error) {
	sa, err := ResolveAddr(addr)
	if err != nil {
		return nil, errors.New("resolve addr err: " + err.Error())
	}
	switch proto {
	case "udp":
		return new(UDP).Dial(sa)
	default:
		return nil, errors.New("undefined proto " + proto)
	}
}
