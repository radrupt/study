统计当天登录用户数量
通过计算用户id的hash结果，
```
package main

import "fmt"
import "hash/adler32"
func CheckECC(data []byte) uint32 {
	if len(data) > 200 {
		tmpdata := append(data[:100], data[len(data)-100:]...)
		return adler32.Checksum(tmpdata)
	} else {
		return adler32.Checksum(data)
	}
}
func main() {
	uidHash := CheckECC([]byte("23")) // 获取用户对应的唯一数字
	// 将该数字作为offset放入bitmap
	// 之后通过bitcount来获取uv
}
// 第一个参数就是offset，第二个参数设为1表示有用户
> SETBIT mykey 7 1
> SETBIT mykey 2 1
> bitcount mykey
// 2
```