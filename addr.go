package halo

import (
	"fmt"

	"golang.org/x/sys/unix"
)

type Addr struct {
	IP   [4]byte
	Port int
}

func (a *Addr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d:%d", a.IP[0], a.IP[1], a.IP[2], a.IP[3], a.Port)
}

func (a *Addr) IPString() string {
	return fmt.Sprintf("%d.%d.%d.%d", a.IP[0], a.IP[1], a.IP[2], a.IP[3])
}

func (a *Addr) ToSockaddrInet4() *unix.SockaddrInet4 {
	return &unix.SockaddrInet4{
		Port: a.Port,
		Addr: a.IP,
	}
}

func (a *Addr) ParseSockaddrInet4(sa *unix.SockaddrInet4) *Addr {
	a.IP = sa.Addr
	a.Port = sa.Port
	return a
}
