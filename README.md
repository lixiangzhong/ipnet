# ipnet


```go
package main

import (
	"fmt"
	"github.com/lixiangzhong/ipnet"
)

func main() {
	ip, err := ipnet.ParseIPv4("1.1.1.1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ip, ip.Int(), ip.Inverse())
	ip.ParseInt(111)
	fmt.Println(ip)
	cidr := ipnet.MustParseCIDR("1.1.1.0/24")
	fmt.Println(cidr)
	fmt.Println(cidr.Int())
	fmt.Println(cidr.IPMask())
	fmt.Println(cidr.StartEndIP())
	cidr, err = ipnet.IPMaskToCIDR("1.1.0.0", "255.255.0.0")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cidr)
	fmt.Println(cidr.IPMask())
	cidrs, err := ipnet.IPRangeToCIDR("1.1.1.0", "1.1.2.255")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cidrs)
}
```