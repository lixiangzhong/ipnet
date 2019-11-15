package ipnet

import (
	"encoding/binary"
	"math/bits"
	"net"
)

type IPv4 struct {
	net.IP
	byteorder binary.ByteOrder
}

func MustParseIPv4(s string) IPv4 {
	ip, err := ParseIPv4(s)
	if err != nil {
		panic(err)
	}
	return ip
}

func ParseIPv4(s string) (IPv4, error) {
	var ip IPv4
	err := ip.Parse(s)
	return ip, err
}

func ParseIPv4FromUint32(i uint32) IPv4 {
	var ip IPv4
	ip.ParseInt(i)
	return ip
}

func (i *IPv4) SetByteOrder(b binary.ByteOrder) {
	i.byteorder = b
}

func (i IPv4) Int() uint32 {
	if i.byteorder == nil {
		return binary.BigEndian.Uint32(i.IP.To4())
	}
	return i.byteorder.Uint32(i.IP.To4())
}

func (i *IPv4) Parse(s string) error {
	i.IP = net.ParseIP(s).To4()
	if i.IP == nil {
		return ErrIPv4Format(s)
	}
	return nil
}

func (i *IPv4) ParseInt(u uint32) {
	var ip = make(net.IP, 4)
	if i.byteorder == nil {
		binary.BigEndian.PutUint32(ip, u)
	} else {
		i.byteorder.PutUint32(ip, u)
	}
	i.IP = ip
}

func (i *IPv4) ParseBytes(b []byte) {
	if i.byteorder == nil {
		i.ParseInt(binary.BigEndian.Uint32(b))
	} else {
		i.ParseInt(i.byteorder.Uint32(b))
	}
}

func (i IPv4) Equal(x *IPv4) bool {
	return i.IP.To4().Equal(x.IP.To4())
}

func (i *IPv4) Set(a, b, c, d byte) {
	ip := net.IPv4(a, b, c, d)
	i.IP = ip.To4()
}

func (i *IPv4) SetA(a byte) {
	ip := i.IP.To4()
	ip[0] = a
	i.IP = ip
}

func (i *IPv4) SetB(b byte) {
	ip := i.IP.To4()
	ip[1] = b
	i.IP = ip
}

func (i *IPv4) SetC(c byte) {
	ip := i.IP.To4()
	ip[2] = c
	i.IP = ip
}

func (i *IPv4) SetD(d byte) {
	ip := i.IP.To4()
	ip[3] = d
	i.IP = ip
}

func (i IPv4) Inverse() IPv4 {
	var ip IPv4
	ip.byteorder = i.byteorder
	ip.ParseInt(^i.Int())
	return ip
}

func (i IPv4) AddInt(n uint32) IPv4 {
	return ParseIPv4FromUint32(i.Int() + n)
}

func (i IPv4) SubInt(n uint32) IPv4 {
	return ParseIPv4FromUint32(i.Int() - n)
}

func (i IPv4) Ones() int {
	return 32 - bits.TrailingZeros32(i.Int())
}

type ErrIPv4Format string

func (e ErrIPv4Format) Error() string {
	return "incorrect ipv4 format " + string(e)
}
