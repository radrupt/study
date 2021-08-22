1. atomic.LoadInt32(addr *int32) 能原子的读数据，即在读期间是不允许其它cpu对数据进行写操作的
2. atomic.StoreInt32(addr *int32, val int32) 能原子的存数据，即在存期间，不允许其它cpu对数据进行读写操作
3. atomic.CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
4. atomic.Value
```
var v atomic.Value
v.Store(100)
v.Load(100).(int32)
```
5. atomic.AddInt32(addr *int32, delta int32) (val int32)