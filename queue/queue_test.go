package queue

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var wg sync.WaitGroup

type Data struct {
	id  int
	msg string
}

func init() {
	wg = sync.WaitGroup{}
}

func TestPop(t *testing.T) {
	queue := New(1e5) // 10w

	for i := 0; i < 5; i++ { // 超出不存
		d := Data{id: i, msg: "hello world"}
		queue.Push(d)
	}
	fmt.Println("1-->", queue.Len())
	for i := 0; i < 3; i++ {
		fmt.Println("1-->", queue.Pop())
	}

	fmt.Println("2-->", queue.Len())

}

func TestBasic(t *testing.T) {
	queue := New(1e5) // 10w

	for i := 0; i < 1e4; i++ { // 超出不存
		d := Data{id: i, msg: "hello world"}
		queue.Push(d)
	}
	fmt.Println("1-->", queue.Len())
	for i := 0; i < 1e4; i++ {
		fmt.Println("1-->", queue.Pop())
	}

	fmt.Println("2-->", queue.Len())

}

func TestQueue(t *testing.T) {
	Convey("queue", t, func() {
		limitQueue := New(4)
		// 一般 push
		So(limitQueue.Push(1), ShouldBeTrue)
		So(limitQueue.PushBatch([]interface{}{2, 3, 4}), ShouldBeTrue)
		// 超出容量
		So(limitQueue.Push(5), ShouldBeFalse)
		So(limitQueue.PushBatch([]interface{}{6, 7}), ShouldBeFalse)
		// 读出数据
		So(limitQueue.Len(), ShouldEqual, 4)
		So(limitQueue.Pop(), ShouldEqual, 1)
		So(limitQueue.Len(), ShouldEqual, 3)
		// batch 读出
		So(limitQueue.PopBatch(2), ShouldResemble, []interface{}{2, 3})
		So(limitQueue.PopBatch(100), ShouldHaveLength, 1)
		So(limitQueue.PopBatch(100), ShouldHaveLength, 0)
	})

	Convey("dumper", t, func() {
		dumperQueue := NewDumper(5, "./xxx")
		dumperQueue.SetDir("./testdata")
		So(dumperQueue.PushBatch([]interface{}{1, 2, 3, 4, 5, 6}), ShouldBeTrue)
		// 使用 0，全量输出
		dumpSize, err := dumperQueue.Dump(0)
		So(err, ShouldBeNil)
		So(dumpSize, ShouldEqual, 6)

		// 已经没有数据，dumpSize 为 0
		dumpSize, err = dumperQueue.Dump(0)
		So(err, ShouldBeNil)
		So(dumpSize, ShouldEqual, 0)

		// 读取回来，对比数据
		err = dumperQueue.ReadOnce(func(q Queue, data []byte) error {
			So(bytes.Equal(data, []byte("[1,2,3,4,5,6]")), ShouldBeTrue)
			return nil
		})
		So(err, ShouldBeNil)

		os.RemoveAll("./xxx")
		os.RemoveAll("./testdata")
	})
	Convey("chan-queue", t, func() {
		cQueue := NewChan() // 容量是有用的
		// 一般 push
		So(cQueue.Push(1), ShouldBeTrue)
		So(cQueue.PushBatch([]interface{}{2, 3, 4}), ShouldBeTrue)
		So(cQueue.Push(5), ShouldBeTrue)
		So(cQueue.PushBatch([]interface{}{6, 7}), ShouldBeTrue)
		close(cQueue.Ch)
		// 读出数据
		So(cQueue.Len(), ShouldEqual, 7)
		So(cQueue.Pop(), ShouldEqual, 1)
		So(cQueue.Len(), ShouldEqual, 6)
		// batch 读出
		So(cQueue.PopBatch(2), ShouldResemble, []interface{}{2, 3})
		So(cQueue.PopBatch(100), ShouldHaveLength, 4)
		So(cQueue.PopBatch(100), ShouldHaveLength, 0)
	})
}
