package halo

import (
	"golang.org/x/sys/unix"
)

type TCP struct{ Fd int }

func (t *TCP) Listen(sa *Addr) (SConn, error) {
	if fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, unix.IPPROTO_TCP); err != nil {
		return nil, err
	} else {
		defer unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
		if err := unix.Bind(fd, sa.ToSockaddrInet4()); err != nil {
			return nil, err
		}
		if err := unix.Listen(fd, 128); err != nil {
			return nil, err
		}
		t.Fd = fd
		return t, nil
	}
}

func (t *TCP) Accept() (Conn, error) {
	if fd, _, err := unix.Accept(t.Fd); err != nil {
		return nil, err
	} else {
		return &TCPConn{Fd: fd}, nil
	}
}

func (t *TCP) Dial(sa *Addr) (Conn, error) {
	if fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, unix.IPPROTO_TCP); err != nil {
		return nil, err
	} else {
		if err = unix.Connect(fd, sa.ToSockaddrInet4()); err != nil {
			return nil, err
		}
		return &TCPConn{Fd: fd}, nil
	}
}

func (t *TCP) Close() error {
	return unix.Close(t.Fd)
}

type TCPConn struct{ Fd int }

func (conn *TCPConn) Read(b []byte) (n int, err error) {
	return unix.Read(conn.Fd, b)
}

func (conn *TCPConn) Write(b []byte) (n int, err error) {
	return unix.Write(conn.Fd, b)
}

func (conn *TCPConn) Close() error {
	return unix.Close(conn.Fd)
}
