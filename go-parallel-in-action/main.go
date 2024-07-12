package main

import (
	"fmt"
	"time"
)

// func main() {
// 	a := struct {
// 		Name string
// 		Code string
// 	}{
// 		"laowang",
// 		"golang",
// 	}

// 	// ---- 通用
// 	// 值的默认格式表示
// 	fmt.Printf("%v\n", a)
// 	// 类似%v，但输出结构体时会添加字段名
// 	fmt.Printf("%+v\n", a)
// 	// 值的Go语法表示, 源码片段
// 	fmt.Printf("%#v\n", a)
// 	// 值的类型的Go语法表示
// 	fmt.Printf("%T\n", a)
// 	// 百分号
// 	fmt.Printf("%%\n")

// 	// --- 布尔
// 	// %t    单词true或false
// 	fmt.Printf("%t\n", true)

// 	// --- 整数
// 	// %b    表示为二进制
// 	fmt.Printf("%b\n", 4) // --> 100
// 	// %c    该值对应的unicode码值, 打印字符，如 33 标识!
// 	fmt.Printf("%c\n", 33)
// 	// %q    该值对应的单引号括起来的go语法字符串字面值，必要时会采用安全的转义表示
// 	fmt.Printf("%q\n", 1000)
// 	// %d    表示为十进制
// 	fmt.Printf("%d\n", 1)
// 	// %o    表示为八进制
// 	fmt.Printf("%o\n", 9)
// 	// %x    表示为十六进制，使用a-f
// 	fmt.Printf("%x\n", 15)
// 	// %X    表示为十六进制，使用A-F
// 	fmt.Printf("%X\n", 15)
// 	// %U    表示为Unicode格式：U+1234，等价于"U+%04X"
// 	fmt.Printf("%U\n", 33)

// 	// --- 浮点数
// 	// %b    无小数部分、二进制指数的科学计数法，如-123456p-78；参见strconv.FormatFloat
// 	fmt.Printf("%b\n", 1234.1234)
// 	// %e    科学计数法，如-1234.456e+78
// 	fmt.Printf("%e\n", 1000000.00)
// 	// %E    科学计数法，如-1234.456E+78
// 	fmt.Printf("%E\n", 1000000.00)
// 	// %f    有小数部分但无指数部分，如123.456
// 	fmt.Printf("%f\n", 1234.1234333)
// 	// %F    等价于%f
// 	fmt.Printf("%F\n", 1234.1234333)
// 	// %g    根据实际情况采用%e或%f格式（以获得更简洁、准确的输出）
// 	fmt.Printf("%g\n", 1234.1234333)
// 	// %G    根据实际情况采用%E或%F格式（以获得更简洁、准确的输出）
// 	fmt.Printf("%G\n", 1234.1234333)

// 	// --- 宽度、精度
// 	/*
// 		对于大多数类型的值，宽度是输出字符数目的最小数量，如果必要会用空格填充。对于字符串，精度是输出字符数目的最大数量，如果必要会截断字符串。

// 		对于整数，宽度和精度都设置输出总长度。采用精度时表示右对齐并用0填充，而宽度默认表示用空格填充。

// 		对于浮点数，宽度设置输出总长度；精度设置小数部分长度（如果有的话），除了%g和%G，此时精度设置总的数字个数。例如，对数字123.45，格式%6.2f 输出123.45；格式%.4g输出123.5。%e和%f的默认精度是6，%g的默认精度是可以将该值区分出来需要的最小数字个数。

// 	*/
// 	// %f:    默认宽度，默认精度
// 	// %9f    宽度9，默认精度
// 	fmt.Printf("%9f\n", 12345.1234567890)
// 	// %.2f   默认宽度，精度2
// 	fmt.Printf("%.2ff\n", 12345.1234567890)
// 	//%9.2f  宽度9，精度2
// 	fmt.Printf("%9.2f\n", 1234567890.1234567890)
// 	//%9.f   宽度9，精度0
// 	fmt.Printf("%9.f\n", 12345.1234567890)

// 	// 字符串宽度
// 	fmt.Printf("%6s\n", "123456")
// 	// 用0补齐不足的位
// 	fmt.Printf("%06s\n", "12345")
// 	fmt.Printf("%6.6s\n", "12345678")

// 	// --- 字符串
// 	// %s    直接输出字符串或者[]byte
// 	fmt.Printf("%s\n", "s")
// 	// %q    该值对应的双引号括起来的go语法字符串字面值，必要时会采用安全的转义表示
// 	fmt.Printf("%q\n", "s")
// 	// %x    每个字节用两字符十六进制数表示（使用a-f）
// 	fmt.Printf("%x\n", "a")
// 	// %X    每个字节用两字符十六进制数表示（使用A-F）
// 	fmt.Printf("%X\n", "a")

// 	// --- 指针
// 	// %p    表示为十六进制指针地址，并加上前导的0x
// 	fmt.Printf("%p\n", &a)

// 	// 字符默认右对齐
// 	fmt.Printf("|%6s|%6s|\n", "foo", "b")
// 	// 字符左对齐
// 	fmt.Printf("|%-6s|%-6s|\n", "foo", "b")
// }

// 抽象一个栅栏
type Barrier interface {
	Wait()
}

// 创建栅栏对象
func NewBarrier(n int) Barrier {
	b := barrier{chCount: make(chan struct{}), n: n, chSync: make(chan struct{})}
	go b.Sync()
	return b
}

// 栅栏的实现类
type barrier struct {
	chCount chan struct{}
	chSync  chan struct{}
	n       int
}

// 测试代码
func (b barrier) Sync() {
	count := 0
	for range b.chCount {
		count++
		if count >= b.n {
			fmt.Println("统计结束")
			close(b.chSync)
			break
		}
	}
}

func (b barrier) Wait() {
	b.chCount <- struct{}{}
	<-b.chSync // 阻塞同步器
}

func main() {
	b := NewBarrier(10)
	fmt.Println("开始")
	for i := 0; i < 10; i++ {
		go b.Wait()
	}

	// 模拟常驻
	time.Sleep(time.Second)
}
