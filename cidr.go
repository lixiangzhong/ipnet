package ipnet

import (
	"fmt"
	"math/bits"
	"net"
	"unicode"
)

type ErrCIDRFormat string

func (e ErrCIDRFormat) Error() string {
	return "incorrect CIDR format " + string(e)
}

type CIDR struct {
	*net.IPNet
}

func ParseCIDR(s string) (*CIDR, error) {
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}
	if !ipnet.IP.Equal(ip) {
		return nil, ErrCIDRFormat(s)
	}
	return &CIDR{ipnet}, nil
}

func IPMaskToCIDR(ip string, mask string) *CIDR {
	var ipnet = new(net.IPNet)
	ipnet.IP = net.ParseIP(ip).To4()
	ipnet.Mask = net.IPMask(net.ParseIP(mask).To4())
	return &CIDR{ipnet}
}

func IPRangeToCIDR(startip, endip string) ([]*CIDR, error) {
	start := new(IPv4)
	end := new(IPv4)
	err := start.Parse(startip)
	if err != nil {
		return nil, err
	}
	err = end.Parse(endip)
	if err != nil {
		return nil, err
	}
	endint := end.Int()
	startint := start.Int()

	var cidrs = make([]*CIDR, 0)
	var i int
	for endint >= startint {
		bit := uint32(bits.TrailingZeros32(^endint))
		if i == 0 && bit == 0 || endint == startint {
			ip := new(IPv4)
			ip.ParseInt(endint)
			_, ipnet, _ := net.ParseCIDR(fmt.Sprintf("%s/32", ip))
			cidrs = append(cidrs, &CIDR{ipnet})
			endint--
		}
		i++
		for bit > 0 {
			begin := (endint >> bit) << bit
			if begin < startint {
				bit--
			} else {
				ip := new(IPv4)
				ip.ParseInt(begin)
				_, ipnet, err := net.ParseCIDR(fmt.Sprintf("%s/%v", ip, 32-bit))
				if err != nil {
					return nil, err
				}
				cidrs = append(cidrs, &CIDR{ipnet})
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

func (c *CIDR) SplitTo(tomask int) []*CIDR {
	var cidrs = make([]*CIDR, 0)
	mask, bits := c.IPNet.Mask.Size()
	if tomask <= mask {
		cidrs = append(cidrs, c)
		return cidrs
	}
	ip := &IPv4{IP: c.IP}
	for i := 0; i <= (1<<uint(tomask-mask))-1; i++ {
		n := ip.Int() + uint32(i<<uint(bits-tomask))
		network := new(IPv4)
		network.ParseInt(n)
		_, ipnet, _ := net.ParseCIDR(fmt.Sprintf("%s/%v", network, tomask))

		cidrs = append(cidrs, &CIDR{ipnet})
	}
	return cidrs
}

func (c *CIDR) Int() (uint32, uint32) {
	start, end := c.StartEndIP()
	return start.Int(), end.Int()
}

func (c *CIDR) StartEndIP() (*IPv4, *IPv4) {
	mask, bit := c.Mask.Size()
	startip := new(IPv4)
	startip.Parse(c.IP.String())
	endip := new(IPv4)
	endip.ParseInt(1<<(uint32(bit-mask)) + startip.Int() - 1)
	return startip, endip
}

func (c *CIDR) IPMask() (*IPv4, *IPv4) {
	ip := new(IPv4)
	ip.IP = c.IP
	mask := new(IPv4)
	mask.ParseBytes([]byte(c.IPNet.Mask))
	return ip, mask
}

func isNumber(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) == false {
			return false
		}
	}
	return true
}
