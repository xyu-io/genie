## queue
主要基于`container/list`实现了队列，并添加了一些诸如`批量Push`和`批量Pop`的功能，并添加了读写锁

```
.
# 实现Queue接口，用channel来模拟队列
├── chan.go
# 支持dump的队列实现
├── dumper.go
# Queue的空实现 
├── empty.go
# 基于container/list的封装，定义了Queue接口，实现了baseQueue（带读写锁）和Limited（baseQueue的基础上增加最大数量限制）
├── queue.go
├── queue_test.go
└── readme.md
```