在 `gorm.DB` 的 `Where` 方法中，可以使用多个参数来指定多个筛选条件。这些参数会被解释为 SQL 查询中的多个条件，并使用 `AND` 连接它们。下面是一个示例代码：

```go
db.Where("name = ?", "Alice").Where("age >= ?", 18).Find(&users)
复制代码复制代码复制代码复制代码
```

在这个例子中，`db` 是一个 `gorm.DB` 对象，`Find` 方法会查找 `users` 切片中符合条件的记录。`Where` 方法被调用两次，每次传递了一个条件，这两个条件会被连接起来使用。

你也可以将多个条件放在一个 `Where` 方法中，使用占位符来代替参数，例如：

```go
db.Where("name = ? AND age >= ?", "Alice", 18).Find(&users)
复制代码复制代码复制代码复制代码
```

在这个例子中，两个条件使用 `AND` 连接，都被包含在同一个 `Where` 方法中。

### 软删除与硬删除

1. gorm默认执行软删除。对于内嵌了gorm.DB类型的结构体变量，建立orm后删除时并不会从数据库中把数据删除，而是会设定一个deleteat时间

2. gorm查询的时候默认启动了软删除筛选，只使用Where和Find方法无法检索到被软删除的记录

3. 可以使用UnScoped方法来去除软删除筛选，

   比如`db.UnScoped().Where("deleteat is not null",id).Find(targets)` 

   能够检索到所有被软删除的数据。

   又比如调用了UnScoped()之后再执行Delete就是硬删除