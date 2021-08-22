golang中的类型转换

1. 断言
```
package main

import (
	"fmt"
)
type A []string
func main() {
	var a A
	func (o interface{}) {
		switch o.(type) {
			case A:
				fmt.Println(o.(A))
		}
	}(a)
}
```
2. 强制类型转换

3. 数字转string
如果使用string(120),将无法得到想要的结果，反而做了隐士转换，即120变成了unicode编码的的byte，byte变为string
但其实是直接变成了utf-8知识这里120对应了unicode

正确的用法应该是strconv.Itoa()

