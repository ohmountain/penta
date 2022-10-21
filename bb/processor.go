package bb

import "context"

// WriteReceipt
// 写入数据后返回的回执，它必须包含一个唯一的ID和写入时间戳
type WriteReceipt interface {

	// Id
	// 获取此次写入数据回执的ID
	Id() string

	// Timestamp
	// 获取此次写入数据的时间戳
	Timestamp() int64
}

// OutFunc 数据输出函数
// 这个函数约定了接收已处理数据的函数
type OutFunc[T any] func(out T)

// Processor 数据处理器
// 任意类型数据处理器的约定
// 实现这个接口，需要约定两个类型
// In 输入数据的类型
// Out 输出数据的类型
type Processor[In any, Out any] interface {

	// Id
	// 每一个处理器的实现都必须有一个唯一的Id
	Id() string

	// Write
	// 处理器接收外部数据的入口
	Write(data In) WriteReceipt

	// Run
	// 执行数据处理，这个函数应该使用协程实现
	// 这个函数接收一个context作为停止信号
	Run(ctx context.Context)

	// Exit
	// 退出数据处理
	// 除了通过Run函数的ctx退出，也可以通过调用Exit函数主动退出数据处理
	Exit() error

	// Read
	// 外部接收已处理数据的入口
	// 这个函数需要一个类型为[Out]的参数
	Read(reader OutFunc[Out])
}
