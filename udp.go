package halo

import (
	"golang.org/x/sys/unix"
)

type UDP struct {
	Fd int
	Ra *unix.SockaddrInet4
}

func (u *UDP) Listen(sa *Addr) (PConn, error) {
	if fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, unix.IPPROTO_UDP); err != nil {
		return nil, err
	} else {
		defer unix.SetNonblock(fd, true)
		defer unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
		if err := unix.Bind(fd, sa.ToSockaddrInet4()); err != nil {
			return nil, err
		}
		return &UDP{Fd: fd}, nil
	}
}

func (u *UDP) Dial(sa *Addr) (PConn, error) {
	if fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0); err != nil {
		return nil, err
	} else {
		defer unix.SetNonblock(fd, true)
		return &UDP{Fd: fd, Ra: sa.ToSockaddrInet4()}, nil
	}
}

func (u *UDP) SendMsg(b []byte) error {
	return unix.Sendto(u.Fd, b, 0, u.Ra)
}

func (u *UDP) SendMsgTo(b []byte, sa *Addr) error {
	return unix.Sendto(u.Fd, b, 0, sa.ToSockaddrInet4())
}

func (u *UDP) RecvMsg(b []byte) (n int, sa *Addr, err error) {
	n, from, err := unix.Recvfrom(u.Fd, b, 0)
	return n, new(Addr).ParseSockaddrInet4(from.(*unix.SockaddrInet4)), err
}

func (u *UDP) Close() error {
	return unix.Close(u.Fd)
}
