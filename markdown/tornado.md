##Tornado
直接读代码
### 简介
FriendFeed使用了一款使用Python编写的，相对简单的非阻塞式Web服务器。它的特点在于非阻塞服务器和epoll的运用，这个Web框架还包含了一些相关的有用工具和优化。

以下是经典的“Hello, world”示例：

```

import tornado.ioloop
import tornado.web

class MainHandler(tornado.web.RequestHandler):
    def get(self):
        self.write("Hello, world")

application = tornado.web.Application([
    (r"/", MainHandler),
])

if __name__ == "__main__":
    application.listen(8888)
    tornado.ioloop.IOLoop.instance().start()

```
Tornado各模块的相互依存关系较少，所以理论上讲，可以在自己的项目中独立地使用任何模块，而不需要使用整个包。

###主要模块
###### 核心Web框架

* tornado.web - RequestHandler 和 Application 类
* tornado.httpserver - 非阻塞式HTTP服务器
* tornado.template - 灵活的输出生成器
* tornado.escape - 转义和字符串处理
* tornado.locale - 国际化支持

###### 异步网络
* tornado.ioloop - 核心的事件循环
* tornado.iostream - 对非阻塞式的socket的简单封装，以方便常用读写操作
* tornado.httpclient - 非阻塞式HTTP客户端
* tornado.netutil - 各种网络工具

###### 其他服务的集成
* tornado.auth - OpenID和OAuth的第三方认证登录
* tornado.database - 简单的MySQL客户端封装
* tornado.platform.twisted - Tornado上的Twisted运行代码
* tornado.websocket - 和浏览器的双向通信
* tornado.wsgi - 和其他Python框架和服务器的协作

###### 工具
* tornado.autoreload - 在开发中自动检测代码更改
* tornado.gen - 简单的异步代码
* tornado.httputil - 处理HTTP头和URL
* tornado.options - 命令行解析
* tornado.process - 多进程工具
* tornado.stack_context - 异步回调中的异常处理
* tornado.testing - 异步代码的单元测试支持

具体使用场景和基本操作，请查看文档[tornado中文文档](http://www.tornadoweb.cn/documentation)

### 深入理解Tornado

要读哪些内容？
你可以参看这篇文章来了解需要读哪些文件：Tornado源码必须要读的几个核心文件

Core web framework 部分，tornado.web 包含web框架的大部分主要功能，这个需要重点看，它包含RequestHandler和Application两个重要的类。Application 是个单例，总揽全局路由，创建服务器负责监听，并把服务器传回来的请求进行转发（__call__）。RequestHandler 是个功能很丰富的类，基本上 web 开发需要的它都具备了，比如redirect，flush，close，header，cookie，render（模板），xsrf，etag等等。关于这个模块的解析，可以参看以下三篇文章：

Tornado RequestHandler和Application类
Application对象的接口与起到的作用
RequestHandler的分析
接下来是 tornado.httpserver，一个无阻塞HTTP服务器的实现。从 web 跟踪到 httpserver.py 和 tcpserver.py。这两个文件主要是实现 http 协议，解析 header 和 body， 生成request，回调给 appliaction，一个经典意义上的 http 服务器（written in python）。众所周知，这是个很考究性能的一块（IO），所以它和其它很多块都连接到了一起，比如 IOLoop，IOStream，HTTPConnection 等等。这里 HTTPConnection 是实现了 http 协议的部分，它关注 Connection 嘛，这是 http 才有的。至于监听端口，IO事件，读写缓冲区，建立连接之类都是在它的下层–tcp里需要考虑的，所以，tcpserver 才是和它们打交道的地方。关键部分是 HTTP 层与 TCP 层的工作原理：

HTTP层：HTTPRequest,HTTPServer与HTTPConnection
Tornado在TCP层里的工作机制
剩下的 tornado.template，tornado.escape 之类的可以先不阅读，会使用就行。

接下来是 Asynchronous networking 底层模块，特别是底层的网络模块。

tornado.ioloop 是核心的I/O循环，需要重点看。如果你用过 select/poll/epoll/libevent 的话，对它的处理模型应该相当熟悉。简言之，就是一个大大的循环，循环里等待事件，然后处理事件。这是开发高性能服务器的常见模型，tornado 的异步能力就是在这个类里得到保证的：

Tornado高性能的秘密：ioloop对象分析
Tornado IOLoop instance()方法的讲解
Tornado IOLoop start()里的核心调度
Tornado IOLoop与Configurable类
然后是tornado.iostream，对非阻塞式的 socket 的简单封装，以方便常用读写操作。这个也是重要模块。IOStream。顾名思义，就是负责IO的。说到IO，就得提缓冲区和IO事件。缓冲区的处理都在它自个儿类里，IO事件的异步处理就要靠 IOLoop 了：

对socket封装的IOStream机制概览
IOStream实现读写的一些细节
如果有兴趣，可以阅读更多源码，比如epoll.py。其实这个文件也没干啥，就是声明了一下服务器使用 epoll。选择 select/poll/epoll/kqueue 其中的一种作为事件分发模型，是在 tornado 里自动根据操作系统的类型而做的选择，所以这几种接口是一样的（当然效率不一样）：

预备知识：我读过的对epoll最好的讲解
epoll与select/poll性能，CPU/内存开销对比
 
 
##### 参考资料
- [Tornado官方文档](http://tornadoweb.org)
- [tornado中文文档](http://www.tornadoweb.cn/documentation)
- [深入理解Tornado-一个异步web服务器](http://golubenco.org/understanding-the-code-inside-tornado-the-asynchronous-web-server-powering-friendfeed.html)
