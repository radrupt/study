## go 面向对象编程
> 面向对象三大特点：封装，继承，多态
1. 封装是指将现实世界中存在的某个客体的属性与行为绑定在一起，并放置在一个逻辑单元内，go的interface（接口不是类型）实现了类似的封装
2. 继承简单地说就是一种层次模型，这种层次模型能够被重用。层次结构的上层具有通用性，但是下层结构则具有特殊性。go的匿名字段实现了类似的效果，但实际是组合，因为没有覆盖啊，可在运行时根据上下文即A.Hello()或A.B.Hello()来执行对应的方法。
3. 多态是指不同事物具有不同表现形式的能力。多态机制使具有不同内部结构的对象可以共享相同的外部接口，通过这种方式减少代码的复杂度。go的interface（接口不是类型）算是一种多态
### 面向对象
> go 语言面向对象与 C++，Java 等语言不同之处， 在于 go 不支持继承。Go 语言只支持聚合和嵌入， 通过一个简单的例子来了解一下 go 中的聚合和嵌入,但其实有一种类似继承的形式就是struct A中将struct B作为匿名字段使用，那么在A.Func里使用B.Func可实现调用基类的效果

```
type ColoredPoint struct {
            color.Color // 匿名字段（嵌入）
            x, y int    // 具名字段 (聚合)
} 
```
### 自定义类型
```
type newInt int
type people struct{}
```

### 添加方法
```
type People struct {
    age int
    name string
}
func (this *People) printAge() {
    fmt.Println(this.age)
}

func(this *People) printName() {
    fmt.Println(this.name)
}

func main() {
    jack := People{name:"jack", age:12}
    jack.printAge()
    jack.printName()
}
```

### 重写方法
```
type Fruit struct {
        price int
        quantity int
}

func (fruit *Fruit) Cost() int {
    return fruit.price * fruit.quantity
}

type SpecialFruit struct {
        Fruit  // 匿名字段 （嵌入）
        markup int  // 具名字段 （聚合）
}

// 重写
func (specialFruit *SpecialFruit) Cost() int {
    return specialFruit.Fruit.Cost()  * specialFruit.markup
}

func main() {
    res := SpecialFruit{Fruit{price:12, quantity:12}, 10}
    fmt.Println(res.Cost())
}
```

### 接口
一个类如果实现了一个接口的所有函数，那么这个类就实现了这个接口。
```
type MyInterface interface{
	Print()
}

func TestFunc(x MyInterface) {}

type MyStruct struct {}

func (me MyStruct) Print() {}

func main() {
	var me MyStruct
	TestFunc(me)
}
```

### 浅谈继承和组合
>
	面向对象编程讲究的是代码复用，继承和组合都是代码复用的有效方法。组合是将其他类的对象作为成员使用，继承是子类可以使用父类的成员方法。

	引用一个生动的比方：继承是说“我父亲在家里给我帮了很大的忙”，组合是说“我请了个老头在我家里干活”。

	继承

	在继承结构中，父类的内部细节对于子类是可见的。所以我们通常也可以说通过继承的代码复用是一种“白盒式代码复用”。

	优点：

	简单易用，使用语法关键字即可轻易实现。
	易于修改或扩展那些父类被子类复用的实现。
	缺点：

	编译阶段静态决定了层次结构，不能在运行期间进行改变。
	破坏了封装性，由于“白盒”复用，父类的内部细节对于子类而言通常是可见的。
	子类与父类之间紧密耦合，子类依赖于父类的实现，子类缺乏独立性。当父类的实现更改时，子类也不得不会随之更改。
	组合

	组合是通过对现有的对象进行拼装（组合）产生新的、更复杂的功能。因为在对象之间，各自的内部细节是不可见的，所以我们也说这种方式的代码复用是“黑盒式代码复用”。

	优点：

	通过获取指向其它的具有相同类型的对象引用，可以在运行期间动态地定义（对象的）组合。
	“黑盒”复用，被包含对象的内部细节对外是不可见。不破坏封装，整体类与局部类之间松耦合，彼此相对独立。
	整体类对局部类进行包装，封装局部类的接口，提供新的接口，具有较好的可扩展性。
	缺点：

	整体类不能自动获得和局部类同样的接口，比继承实现需要的代码更多。
	不熟悉的代码的话，不易读懂。
	两者的选择

	is-a关系用继承表达，has-a关系用组合表达。继承体现的是一种专门化的概念而组合则是一种组装的概念。

	个人推荐：除非用到向上转型，不然优先考虑组合。


### golang 使用组合的方式实现继承
> golang并非完全面向对象的程序语言，为了实现面向对象的继承这一神奇的功能，golang允许struct间使用匿名引入的方式实现对象属性方法的组合。Go中组合跟继承唯一的不同在于，继承自其他结构体的struct类型可以直接访问父类结构体的字段和方法。
```
type Pet struct {
  name string
}

type Dog struct {
  Pet
  Breed string
}

func (p *Pet) Speak() string {
  return fmt.Sprintf("my name is %v", p.name)
}

func (p *Pet) Name() string {
  return p.name
}

func (d *Dog) Speak() string {
  return fmt.Sprintf("%v and I am a %v", d.Pet.Speak(), d.Breed)
}

func main() {
  d := Dog{Pet: Pet{name: "spot"}, Breed: "pointer"}
  fmt.Println(d.Name()) // d本身没有Name方法，继承了Pet的Name方法
  fmt.Println(d.Speak()) // d本身有Speak方法，在内部访问了父类Pet的Speak方法
}
```
### 不支持多态
Subtyping 在Java中，Dog继承自Pet，那么Dog类型就是Pet子类。这意味着在任何需要调用Pet类型的场景都可以使用Dog类型替换。这种关系称作多态性，但Go的结构体类型不存在这种机制。
```



转自：https://learnku.com/articles/23418/go-object-oriented-programming

转自：https://blog.csdn.net/zaimeiyeshicengjing/article/details/105968971