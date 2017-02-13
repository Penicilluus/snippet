###### 基础类型

源码在：$GOROOT/src/pkg/runtime/runtime.h

```go
/*
 * basic types
 */
typedef signed char             int8;
typedef unsigned char           uint8;
typedef signed short            int16;
typedef unsigned short          uint16;
typedef signed int              int32;
typedef unsigned int            uint32;
typedef signed long long int    int64;
typedef unsigned long long int  uint64;
typedef float                   float32;
typedef double                  float64;

#ifdef _64BIT
typedef uint64          uintptr;
typedef int64           intptr;
typedef int64           intgo; // Go's int
typedef uint64          uintgo; // Go's uint
#else
typedef uint32          uintptr;
typedef int32           intptr;
typedef int32           intgo; // Go's int
typedef uint32          uintgo; // Go's uint
#endif

/*
 * defined types
 */
typedef	uint8			bool;
typedef	uint8			byte;

// 底层类型(uint8) 静态类型(bool)
// uintptr和intptr是无符号和有符号的指针类型
```

###### rune

rune是int32的别名，通常用于表示unicode字符

###### string类型

底层是一个C的struct结构，对string类型的变量初始化，表示对底层结构的初始化

```go
struct String
{
        byte*   str; //字符数组，不用rune类型是因为golang for循环是针对byte的
        intgo   len; //字符数组的长度
};
```

###### slice类型

```go
struct	Slice
{				// must not move anything
	byte*	array;		// actual data 底层数组
	uintgo	len;		// number of elements 实际存放的个数
	uintgo	cap;		// allocated number of elements 总容量
};
// 使用内建函数make初始化
var slice []int32 = make([]int32, 5, 10) // 第一个参数为切片类型，第二个参数为切片大小，第三个可省略默认为切片大小，一般会指定容量，防止频繁分配内存
// 可以由数组直接生成切片
var array = [...]int32{1, 2, 3, 4, 5}
var slice = array[2:4]
//由于切片指向一个底层数组，所以直接修改切片也会影响原切片的内容
var array = [...]int32{1, 2, 3, 4, 5}
var slice = array[2:4]
slice[0]=234 // array {1，2，234，4，5} slice {234，4}

slice = append(slice, 6, 7, 8)
//append的操作令slice重新分配底层数组，所以此时slice的底层数组不再指向前面定义的array。
```

###### 接口类型

```go
//$GOROOT/src/pkg/runtime/type.h
struct Type
{
	uintptr size;
	uint32 hash;
	uint8 _unused;
	uint8 align;
	uint8 fieldAlign;
	uint8 kind;
	Alg *alg;
	void *gc;
	String *string;
	UncommonType *x;
	Type *ptrto;
};

//$GOROOT/src/pkg/runtime/runtime.h
struct Iface
{
	Itab*	tab;
	void*	data;
};
struct Eface
{
	Type*	type;
	void*	data;
};
struct	Itab
{
	InterfaceType*	inter;
	Type*	type;
	Itab*	link;
	int32	bad;
	int32	unused;
	void	(*fun[])(void);
};
```

interface实际上是一个结构体，包括两个成员，一个是指向数据的指针，一个包含了成员的类型信息。Eface是interface{}底层使用的数据结构。

###### map类型

```go
struct Hmap
{
	uintgo  count;
	uint32  flags;
	uint32  hash0;
	uint8   B;
	uint8   keysize;
	uint8   valuesize;
	uint16  bucketsize;

	byte    *buckets;
	byte    *oldbuckets;
	uintptr nevacuate;
};
// map的底层实现是hashmap，使用的引用传值，
var m = make(map[string]int32, 10)
m["hello"] = 123 //会修改原来的map
```

###### 最佳实践

1. Slices, maps, channels, strings, function values, and interface values 实现机制类似指针，所以可以直接传递

struct{}的惯用方法

```go
a := struct{}{}
println(unsafe.Sizeof(a))
// Output: 0 struct{}{}不会占用内存空间

// When implementing a data set:


```

- 当使用map实现set时，用struct{}{}作为map的值
- 当在遍历一个图时，使用map作为数据接口，struct{}{}作为值
- 当实现一个实体不需要任何值，只在意它的方法集时
- 作为channel中的信号量

swap two values

```go
a, b = b, a;
a := 1; b := 2; 
a, b, a = b, a, b;
a = 2 ;
b = 1 ;
// 在交换值时，首先会复制一份临时的值，然后再进行真正的赋值，跟值的顺序无关
```

copy

```go
a := []int{1, 2}
b := []int{3, 4}
check := a
copy(a, b)
fmt.Println(a, b, check)
// Output: [3 4] [3 4] [3 4] 替换了底层数组，所以check也跟着变了

a := []int{1, 2}
b := []int{3, 4}
check := a
a = b
fmt.Println(a, b, check)
// Output: [3 4] [3 4] [1 2] 没有替换底层数组，只是修改了指针的指向，所以check没有变

a := map[string]bool{"A": true, "B": true}
b := make(map[string]bool)
for key, value := range a {
	b[key] = value
}
// copy 一个map最简单的方法

// Following example copies just the description of the map:
a := map[string]bool{"A": true, "B": true}
b := map[string]bool{"C": true, "D": true}
check := a
a = b
fmt.Println(a, b, check)
// Output: map[C:true D:true] map[C:true D:true] map[A:true B:true] 
```

how compare two struct

```go
// compare two structs with the == operator,make sure the they do not contain any slice、map、or function in which case the code will not be compiled.
type Foo struct {
	A int
	B string
	C interface{}
}
a := Foo{A: 1, B: "one", C: "two"}
b := Foo{A: 1, B: "one", C: "two"}
// 比较struct
println(a == b)
// Output: true

type Bar struct {
	A []int
}
a := Bar{A: []int{1}}
b := Bar{A: []int{1}}

println(a == b)
// Output: invalid operation: a == b (struct containing []int cannot be compared)

// 比较slice，map 可以用reflect.DeepEqual()比较
// Both structs and interfaces which contain maps, slices (but not functions) can be compared with the reflect.DeepEqual() function:
var a interface{}
var b interface{}

a = []int{1}
b = []int{1}
println(reflect.DeepEqual(a, b))
// Output: true

a = map[string]string{"A": "B"}
b = map[string]string{"A": "B"}
println(reflect.DeepEqual(a, b))
// Output: true

temp := func() {}
a = temp
b = temp
println(reflect.DeepEqual(a, b))
// Output: false

// For comparing byte slices  比较bytes数组
// bytes.Equal(), bytes.Compare(), and bytes.EqualFold(). 
```

关于String()

```go
type Orange struct {
	Quantity int
}
// 实现String应该使用 实体(o Orange)而不是指针
func (o *Orange) String() string {
	return fmt.Sprintf("%v", o.Quantity)
}
```

###### 内存分配

`基本策略`：

- 每次从操作系统申请一大块内存
- 将申请到的大块内存按照特定大小预先切分成小块，构成链表
- 为对象分配内存时，只需从大小合适的链表提取一个小块
- 回收对象内存，将小块内存重新归还给原链表，以便复用
- 如果闲置内存过多，则尝试归还部分内存给操作系统，降低整体开销

`内存块`：

- span：由多个地址连续的页组成的大块内存，分配器按页数区分大小不同的span
- object将span按特定大小切分成多个小块，每个小块可存储一个对象，object按8字节倍数分为n种，规格化小块内存，优化分配和复用管理策略

分配器初始化时，会构建对照表存储大小和规格的对应关系，包括用来切分的span页数

`分配器`：

- cache：每个运行期工作线程都会绑定一个cache，用于无锁object分配
- central：为所有cache提供分好的后备span资源
- heap：管理闲置span，需要时向操作系统申请新内存

`分配流程`：

1. 计算待分配对象对应规格(size class)。
2. 从 cache.alloc 数组找到规格相同的 span。
3. 从 span.freelist 链表提取可用 object。
4. 如 span.freelist 为空，从 central 获取新 span。
5. 如 central.nonempty 为空，从 heap.free/freelarge 获取，并切分成 object 链表。
6. 如 heap 没有大小合适的闲置 span，向操作系统申请新内存块。

`释放流程`：

1. 将标记为可回收 object 交还给所属 span.freelist。
2. 该 span 被放回 central，可供任意 cache 重新获取使用。
3. 如 span 已收回全部 object，则将其交还给 heap，以便重新切分复用。
4. 定期扫描 heap 里长时间闲置的 span，释放其占用内存。

作为工作线程私有且不被共享的cache是实现高性能无锁分配的核心，而central的作用是在多个cache间提高object利用率，避免内存浪费

###### channel

消息传递：生产者与消费者模型

单向channel：双向channel在作为参数传递时，可以将channel作为单向channel传递，在函数参数中声明即可，在函数内部都是单向channel

同步：主要使用channel作为信号在不同goroutine种传递，make(chan struct{}) 

- 给一个 nil channel 发送数据，造成永远阻塞，造成死锁
- 从一个 nil channel 接收数据，造成永远阻塞，造成死锁
- 给一个已经关闭的 channel 发送数据，引起 panic
- 从一个已经关闭的 channel 接收数据，立即返回一个零值

对于没有buffer的channel，如果只有发送者或者接收者，发送数据和接收数据就会造成死锁

```go
for elem <- range ch {
}
// 当ch被关闭并且缓冲中的数据都被取出后，for循环会退出
```

channel的实现位于[runtime/chan.c](http://golang.org/src/pkg/runtime/chan.c)文件中。每个channel都由一个Hchan结构定义的，这个结构中有两个非常关键的字段就是recvq和sendq。recvq和sendq是两个等待队列，这个两个队列里分别保存的是等待在channel上进行读操作的goroutine和等待在channel上进行写操作的goroutine。如果channel带有缓存则Hchan结构后跟有几个slot，slot数量即为channel缓存大小，无缓存即无slot

读写channel的操作，第一步都是对channel进行加锁，然后查看recvq和sendq的状态，复制数据，修改recvq和sendq中goroutine的运行状态。

###### 接口 interface

接口相当于一个契约，它规定了一个对象所能提供的一组操作

golang不支持完全的面向对象编程，通过组合实现集成，通过接口实现多态

非侵入式接口的优点在于，接口定义与类型定义的分离，

封装：  If a type is only implements the methods of a given interface, it's ok to export only the interface without the underlying type. Obviously, this is helpful in maintaining a cleaner and concise API.

interface{}：it accepts variable number of empty interface(`interface{}`) values.

type assertion：there's a special expression, which let's you assert the type of the value interface holds. This is known as *Type Assertion*.// interface.(*human)，type switches

结构：tab + data， data指向实际引用的数据，tab指向itable，itable首先是描述type信息的一些元数据，然后是满足Stringger接口的函数指针列表（不是实际类型Binary的函数指针集）

http://studygolang.com/articles/2268

如何对goroutine限速，可以构造固定buffer的channel，即可限速

```go
// concurrency-limiting counting semaphore
var sema = make(chan struct{}, 20) // 限速作用
func xxx(){
  	select {
    case sema <- struct{}{}: // acquire token
    case <-done:
        return nil // cancelled
    }
    defer func() { <-sema }() // release token
}


```

###### Golang-反射

- 接口值反射对象
- 反射对象到接口值
- 修改反射对象
- 获取结构体标签
- interface{}到函数反射

参考资料

- [golang 常用数据类型与结构](https://my.oschina.net/goal/blog/196891)
- [CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments)
- [golang编码规范](https://gocn.io/article/1)
- [idiomatic-go](https://dmitri.shuralyov.com/idiomatic-go)
- [深入Go语言网络库的基础实现](http://skoo.me/go/2014/04/21/go-net-core)
- [golang reflect](http://www.cnblogs.com/coder2012/p/4881854.html)