package ipnet

import (
	"fmt"
	"math/bits"
	"net"
)

type ErrCIDRFormat string

func (e ErrCIDRFormat) Error() string {
	return "incorrect CIDR format " + string(e)
}

type CIDR struct {
	*net.IPNet
}

func MustParseCIDR(s string) CIDR {
	c, err := ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return c
}

func ParseCIDR(s string) (CIDR, error) {
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return CIDR{}, err
	}
	if !ipnet.IP.Equal(ip) {
		return CIDR{}, ErrCIDRFormat(s)
	}
	return CIDR{ipnet}, nil
}

func IPMaskToCIDR(ip string, mask string) (CIDR, error) {
	var cidr CIDR
	ipv4, err := ParseIPv4(ip)
	if err != nil {
		return cidr, err
	}
	ipmask, err := ParseIPv4(mask)
	if err != nil {
		return cidr, err
	}
	return ParseCIDR(fmt.Sprintf("%v/%v", ipv4, ipmask.Ones()))
}

func IPRangeToCIDR(startip, endip string) ([]CIDR, error) {
	start, err := ParseIPv4(startip)
	if err != nil {
		return nil, err
	}
	end, err := ParseIPv4(endip)
	if err != nil {
		return nil, err
	}
	endint := end.Int()
	startint := start.Int()

	var cidrs = make([]CIDR, 0)
	var i int
	for endint >= startint {
		bit := uint32(bits.TrailingZeros32(^endint))
		if i == 0 && bit == 0 || endint == startint {
			var ip IPv4
			ip.ParseInt(endint)
			cidr, err := ParseCIDR(fmt.Sprintf("%s/32", ip))
			if err != nil {
				return nil, err
			}
			cidrs = append(cidrs, cidr)
			endint--
		}
		i++
		for bit > 0 {
			begin := (endint >> bit) << bit
			if begin < startint {
				bit--
			} else {
				var ip IPv4
				ip.ParseInt(begin)
				cidr, err := ParseCIDR(fmt.Sprintf("%s/%v", ip, 32-bit))
				if err != nil {
					return nil, err
				}
				cidrs = append(cidrs, cidr)
				if begin == startint {
					return cidrs, nil
				}
				endint = begin - 1
				break
			}

		}
	}
	return cidrs, nil
}

func (c CIDR) SplitTo(tomask int) []CIDR {
	var cidrs = make([]CIDR, 0)
	mask, bits := c.IPNet.Mask.Size()
	if tomask <= mask {
		cidrs = append(cidrs, c)
		return cidrs
	}
	ip := IPv4{IP: c.IP}
	for i := 0; i <= (1<<uint(tomask-mask))-1; i++ {
		n := ip.Int() + uint32(i<<uint(bits-tomask))
		var network IPv4
		network.ParseInt(n)
		cidr := MustParseCIDR(fmt.Sprintf("%s/%v", network, tomask))
		cidrs = append(cidrs, cidr)
	}
	return cidrs
}

func (c CIDR) Int() (uint32, uint32) {
	start, end := c.StartEndIP()
	return start.Int(), end.Int()
}

func (c CIDR) StartEndIP() (IPv4, IPv4) {
	mask, bit := c.Mask.Size()
	var startip IPv4
	startip.IP = c.IP
	var endip IPv4
	endip.ParseInt(1<<(uint32(bit-mask)) + startip.Int() - 1)
	return startip, endip
}

func (c CIDR) IPMask() (IPv4, IPv4) {
	var ip IPv4
	ip.IP = c.IP
	var mask IPv4
	mask.ParseBytes([]byte(c.IPNet.Mask))
	return ip, mask
}
