package ulog

import (
	"math/rand"
	"time"
)

type EventType int8

const (

	// 事件类型
	// 具体事件类型会有一个生命周期,从 BeforeEvent -> ProcessEvent -> AfterEvent
	// BeforeEvent, AfterEvent 很好理解,在打印Info日志的时候会执行前和执行后调用
	// 相对特殊的就是 ProcessEvent, 如果存在两个io.Writer对象,比如同时想写入文件和控制台时
	// 两者等于是执行了两遍Info,所以Info执行过程中的事件会执行两遍
	ProcessEvent EventType = iota // ProcessEvent 执行中的事件
	AfterEvent                    // AfterEvent 执行结束的事件
	BeforeEvent                   // BeforeEvent 执行开始之前的事件
)

func init() {
	events = make([]*event, 0)
}

type event struct {
	id      int64
	etype   EventType
	handler func(*EventType, string)
}

// execute 执行
func (e *event) execute(msg string) {
	e.handler(&e.etype, msg)
	// // 执行结束后要销毁
	// for i, e2 := range events {
	// 	if e2.id == e.id {
	// 		events = append(events[:i], events[i+1])
	// 	}
	// }
}

var events []*event

func addEventMonitor(e *event) {
	events = append(events, e)
}

func getEventMonitor(eType EventType) []*event {
	eMonitors := make([]*event, 0)
	for _, e := range events {
		if e.etype == eType {
			eMonitors = append(eMonitors, e)
		}
	}
	return eMonitors
}

func EventBuilder(t EventType, h func(*EventType, string)) {
	id := genEventId()
	addEventMonitor(&event{id: id, etype: t, handler: h})
}

func genEventId() int64 {
	id := genRand()
	//生成10个0-99之间的随机数
	for i := 0; i < 10; i++ {
		for _, e := range events {
			if e.id == id {
				id = genRand()
			} else {
				return id
			}
		}
	}
	return id
}

func genRand() int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(100000000))
}
