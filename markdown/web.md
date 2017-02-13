golang session 管理

session的创建过程

- 生成全局唯一标识符
- 开辟储存空间，使用redis或者内存空间
- 将唯一标识符发送到客户端（http数据主要可以存储在请求行、头域、body）

将标识符sessionId发送到客户端

- Cookie 服务端通过设置Set-cookie头就可以将session的标识符传送到客户端
- 在返回给用户的页面里的所有的URL后面追加session标识符

参考

- [GO如何实现session](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/06.2.md)



