学习笔记

- Packages that are reusable across many projects only return root error values.
- If the error is not going to be handled, wrap and return up the call stack.
- Once an error is handled, it is not allowed to be passed up the call stack any longer.

如何跑这个应用？

```
go run cmd/main.go
```
启动之后，直接在浏览器访问: http://localhost:8080/user/:id

加了10% random DB Error, 遇到这种error 会打印stacktrace, 如果确认是sql.ErrNoRows则无需打印整个stack trace, 只要打印一下log 是找不到user 即可。 
