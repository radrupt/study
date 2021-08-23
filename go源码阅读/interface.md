1. struct实现接口的两种方式
```
// 通过值接收者实现接口方法
// 通过指针接收者实现接口方法
// 区别：
1. 通过指针接收者实现，可以修改变量，通过值接收者实现不可修改变量
2. 通过指针接收者实现的情况，不可用值类型来调用接口方法，通过值接收者实现的情况，可用值类型或指针类型来调用接口方法
```
2. struct的方法，非实现了接口的方法
由于go语法糖的存在，值类型指针类型都可以调用值接收者实现的方法和指针接收者实现的方法
3. iface和eface的区别
iface支持方法
```
type iface struct {
	tab  *itab
	data unsafe.Pointer
}
type itab struct {
	inter  *interfacetype
	_type  *_type
	link   *itab
	hash   uint32 // copy of _type.hash. Used for type switches.
	bad    bool   // type does not implement interface
	inhash bool   // has this itab been added to hash?
	unused [2]byte
	fun    [1]uintptr // variable sized
}
type interfacetype struct {
	typ     _type
	pkgpath name
	mhdr    []imethod
}
```
eface不支持方法
```
type eface struct {
    _type *_type
    data  unsafe.Pointer
}
```
4. 接口转换的原理
当判定一种类型是否满足某个接口时，Go 使用类型的方法集和接口所需要的方法集进行匹配，如果类型的方法集完全包含接口的方法集，则可认为该类型实现了该接口。
例如某类型有 m 个方法，某接口有 n 个方法，则很容易知道这种判定的时间复杂度为 O(mn)，Go 会对方法集的函数按照函数名的字典序进行排序，所以实际的时间复杂度为 O(m+n)。
5. struct中使用匿名interface
参考：
https://liujiacai.net/blog/2020/03/14/go-struct-interface/