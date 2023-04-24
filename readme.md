## 准备数据库账号

创建账号

`create user 'gomoku'@'%' identified by 'Gomoku:666';`

关闭账号多余访问权限

`revoke all on . from gomoku;`

创建数据库gomokudb

`create database gomokudb;`

指定账号对数据库gomokudb拥有所有访问权限

`grant all on gomokudb.* to gomoku;`

## 创建以及发布前端页面

### 遇到的问题

传递css文件后但是网页却没有加载css样式表， 补充w.Header().Set("Content-Type", "text/css")后解决，因为如果没有设置Content-Type的话浏览器不会把收到的内容当作css解析

#### 资料

Header的Content-Type字段指定了HTTP消息中携带的实体内容类型。常见的Content-Type包括：
text/plain 纯文本格式
text/html HTML格式
text/xml XML格式
image/gif GIF图片格式
image/jpeg JPEG图片格式
image/png PNG图片格式
application/x-www-form-urlencoded 在HTML表单中，把所有数据编码成URL的形式提交到服务器，常见于POST请求
multipart/form-data 在HTML表单中，可以上传文件的编码方式，常见于文件上传
application/json JSON格式

### 前端轮询处理

从polling 到long polling 到websocket

## 后端结构

### 使用orm

#### 问题一.私有属性映射问题

使用go的一个orm包 gorm，但是遇到了一个问题，要表化一个结构体的时候发现私有属性不能够自动表化到映射到数据库的表中，试着`gorm:"column:columnname"`标签修饰私有属性来解决，但是失败了，删去旧的表再度运行的时候私有属性还是没有映射到数据库中。阅读官网文档无果，决定把私有属性改成共有属性

#### 问题二，映射后表中列名字母大小写注意

无论原来字段名是大写还是小写，映射为表中的列名后均为小写

#### 问题三。gorm unsuport &\[\] datatype

有些类型的属性gorm的orm不支持

#### 问题四。

invalid field found for struct cncsmonster/gomoku/model.ChessBoard's field Steps: define a valid foreign key for relations or implement the Valuer/Scanner interface

## 前后端交流

### 关于http请求的收发

1. 注意设置ContentType

2. 设置返回的status,http常见status有:

   ```
   HTTP请求的响应状态码（status code）有很多，以下是一些常见的状态码：
   
    200 OK：请求成功处理，常用于GET、POST请求
    201 Created：请求成功并创建了新资源
   
    204 No Content：请求成功处理，但响应报文中不包含实体的主体部分
    301 Moved Permanently：永久性重定向，请求的资源已被永久性转移到新URI
    302 Found：临时性重定向，请求的资源已被暂时性转移到新URI
    304 Not Modified：缓存资源未修改，服务器允许访问，但未返回新的实体
    400 Bad Request：请求报文中存在语法错误
   
    401 Unauthorized：未授权，需要进行身份验证
   
    403 Forbidden：服务器拒绝该请求，客户端没有访问权限
   
    404 Not Found：服务器上没有找到请求的资源
   
    500 Internal Server Error：服务器在执行请求时遇到错误
   
    503 Service Unavailable：服务器无法处理请求，可能是因为过载或维护
   
    这些状态码对于了解请求处理结果很有用，可以根据状态码来判断请求是否成功，如果不成功，可以根据状态码来了解失败的具体原因。
   ```

3. 

### 约定

1. 使用BoardCase结构保存当前棋盘信息，使用该结构体序列化的json字符串作为browser和server之间传递信息的媒介:

   ```
   type BoardCase struct {
        Curplayer uint    `json:"curplayer"`
        Winner    int     `json:"winner"`
        Chesses   [][]int `json:"chesses"`
    }
   ```

   其中Curplayer属性当前下棋的玩家，可能为:{1,2},
   在web视角，为1表示当前为自己下棋，为2表示当前为敌人下棋。
   在服务端视角,1表示player1下棋,2表示player2下棋

   参考以下资料

   ```
   gorm自动生成的ID通常是一个无符号整型(uint)，并且在默认情况下使用自增主键，从1开始连续分布。在使用MySQL等数据库时，gorm会自动创建一个名为auto_increment的列作为主键，并将其设置为自增。因此，每当您向表中插入新记录时，该列将自动递增并生成一个唯一的ID。此外，gorm还支持使用UUID、雪花算法等自定义主键生成方式。
   ```

   可知用户的ID不会小于1，所以我们可以通过修改ChessBoard中chesses中player1ID的棋子为1，player2ID的棋子为2，小于0的棋子(表示棋盘)为0，来获得player1的BoardCase
   如果要获取player2的BoardCase,就把当前的curPlayer变成对面，把棋盘左右翻转，把1变成2，2变成1

   在浏览器轮询的时候请求获取该信息，使用该信息更新当前棋盘

2. 在浏览器进行下棋操作的时候，使用如下obj对应的json向远程服务器发送消息

   ```
   type Step struct{
       X,
       Y uint
   }
   ```

   服务器解析发来的json。
   如果走棋成功，服务器返回statusOk,
   如果走棋失败，服务器返回status not modified

### 日志

1. fix Playable 4.23
2. 迁移到linux服务器上，通过设置go代理goproxy.cn解决包拉不下的问题
3. 2023.4.24 增加手动切断时存储数据功能
