## 使用net./http

Go的`net/http`包提供了方便的方式来创建和处理HTTP请求和响应。其原理基于Go的`net`包，它提供了底层的网络通信和I/O操作。

当你创建一个HTTP服务器时，你会创建一个`http.Server`对象，它会监听指定的TCP端口，等待连接。当客户端发起连接时，服务器会创建一个新的goroutine去处理这个连接，然后将连接转发到一个`http.ServeMux`对象（也称为HTTP多路复用器）。这个对象负责将请求路由到指定的处理器函数。

每个处理器函数都必须满足以下签名：`func(http.ResponseWriter, *http.Request)`。`http.ResponseWriter`提供了一个接口来写回响应数据，而`*http.Request`包含了请求的信息，如请求方法、头部、URL参数和主体内容等。

在处理请求时，服务器通常会检查请求的方法和路径，并根据它们选择合适的处理器函数。如果找不到合适的处理器函数，服务器会返回一个404错误。

处理器函数可以使用`http.ResponseWriter`对象向客户端发送响应数据。通常，这个对象的`WriteHeader`方法用于发送HTTP状态码和头部信息，而`Write`方法用于发送响应的主体内容。当所有的数据都发送完成时，处理器函数必须调用`http.ResponseWriter`对象的`Close`方法来关闭连接。

如果请求需要发送数据（如表单数据或JSON），则必须在`http.Request`对象中设置请求头部和主体内容。`net/http`包提供了许多辅助函数来简化这些操作，如`http.Post`、`http.Get`、`http.NewRequest`等。

在处理HTTP响应时，服务器必须确保正确设置状态码、头部和主体内容，并在必要时关闭连接。这个过程可以通过使用`http.ResponseWriter`对象的`WriteHeader`、`Header`和`Write`方法来完成。

总之，`net/http`包提供了一个简单而灵活的方式来创建和处理HTTP请求和响应。它提供了一组丰富的工具和函数，可以使开发人员快速构建Web应用程序。