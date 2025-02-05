## lab4
### 需求
- 实现一个key/value存储服务，有多个Raft进行复制
- Get/Put/Append 方法的调用是线性化的，同时保证多个Raft上的执行顺序相同，与lab2不同Put/Append不需要返回值
- 需要处理通过rpc调用raft.Start()后，Leader失去领导权的情况，对于出现分区的情况，可以无限等待直到分区结束。
- 需要实现消息的重复检测

### 通关标准

- 400秒的实际运行时间和600秒的CPU时间（来自lab4的最后一个`Hint`），以及TestSnapshotSize测试小于20秒。
- 所有测试1000次运行稳定通过

有必要在此之前读下助教的文字
https://thesquareplanet.com/blog/students-guide-to-raft/#applying-client-operations