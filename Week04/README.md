学习笔记

Google API Design Guide: [https://cloud.google.com/apis/design/standard_methods]

按照老师介绍的最新的folder 结构，
```
-biz
-data
-service
```
使用了google 的wire 来帮助service 装配

暂时没有用真实的DB 实现Data 层，用了一个map 来模拟，即使真正实现DB做persistence 层， 这个基于map 的实现还可以当作Mock来帮助测试。

没有花时间写test， 只是一个collision 检测逻辑有点复杂所以写了一些happy path 的测试

基于数据库的测试并行一直有些疑问，不知怎么做比较好。

Run Server：
```
go run cmd/serverd/main.go cmd/serverd/wire_gen.go
```

Run Client:
```
go run cmd/cli/main.go
```