# ipnet


```go
package main

import (
	"fmt"
	"github.com/lixiangzhong/ipnet"
)

func main() {
	cidr, err := ipnet.ParseCIDR("1.1.1.0/24")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cidr.Int())
	fmt.Println(cidr.StartEndIP())
	ip, mask := cidr.IPMask()
	fmt.Println(ip, mask, mask.Inverse())
	fmt.Println(mask.Ones())
	fmt.Println(ipnet.IPRangeToCIDR("1.1.1.1", "1.1.1.11"))
}
```