## 学习笔记

- Packages that are reusable across many projects only return root error values.
- If the error is not going to be handled, wrap and return up the call stack.
- Once an error is handled, it is not allowed to be passed up the call stack any longer.

## 如何跑这个应用？

```
go run cmd/main.go
```
启动之后，直接在浏览器访问: http://localhost:8080/user/:id

加了10% random DB Error, 遇到这种error 会打印stacktrace, 如果确认是sql.ErrNoRows则无需打印整个stack trace, 只要打印一下log 是找不到user 即可。

### 测试用例：

1. 400 Response: http://localhost:8080/user/notvalid
1. 404 Response: http://localhost:8080/user/1000
1. 200 Response: http://localhost:8080/user/1   90% 几率
1. 500 Response: http://localhost:8080/user/1   10% 几率
1. Recovery Protection: http://localhost:8080/user/-1


## 作业题目：
我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

答：按照和第三方包交互出现error需要wrap并返回上层处理的原则，应该wrap 并抛给上层。服务层不需要处理，直接抛给最上层的controller 做最后的处理。看了其他一些同学答案，他们觉得要隐藏sql.ErrNoRows，这样保证内部使用的data store 和上层保持透明，觉得也有一定道理，但我觉得不应该在DAO 层吞这个错误，在service 层吞更合理。但是这有点违背只在一处处理错误的原则，还要和老师商讨，暂时在service 层做一下转化，移除掉controller 层对sql 包的依赖。

