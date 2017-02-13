###### Tcp协议

服务器接收到FIN关闭请求后进入CLOSE_WAIT状态并发送ACK返回客户端，进入LAST_ACK状态发送FIN到客户端，当接收到客户端的ACK后进入CLOSE，服务器经历的状态有 CLOSE_WAIT — LAST_ACK — CLOSED

这里最重要的是FIN关闭请求和ACK确认关闭，是客户端和服务器端都需要发送的

客户端经历的状态有 FIN_WAIT1 — FIN_WAIT2 — TIME_WAIT — CLOSE

TIME_WAIT 是客户端收到服务器的关闭请求，关闭本身连接，从TIME_WAIT到CLOSED 一般来说为保证关闭请求正常达到，时间间隔是两个MSL(Max Segment Lifetime)，默认为4分钟，所以在并发量大的时候，如果是服务器端主动关闭连接会造成大量的TIME_WAIT状态。

[Linux Socket过程详细解释（包括三次握手建立连接，四次握手断开连接）](http://www.cnblogs.com/cy568searchx/p/4211124.html)

为什么是三次握手

为了防止已失效的连接请求报文段突然又传送到了服务端，因而产生错误。防止了服务器端的一直等待而浪费资源。

###### Http协议

http中的keep-alive

HTTP/1.0 默认不支持持久连接，很多 HTTP/1.0 的浏览器和服务器使用「Keep-Alive」这个自定义说明来协商持久连接：浏览器在请求头里加上 Connection: Keep-Alive，服务端返回同样的内容，这个连接就会被保持供后续使用。对于 HTTP/1.1，Connection: Keep-Alive 已经失去意义了，因为 HTTP/1.1 除了显式地将 Connection 指定为 close，默认都是持久连接。

跨域问题

`概念定义`：一个资源会发起一个跨域HTTP请求(Cross-site HTTP request), 当它请求的一个资源是从一个与它本身提供的第一个资源的不同的域名时 。域名、协议、端口三者任一不同都会造成前端跨域问题。

跨域并非浏览器限制了发起跨站请求，而是跨站请求可以正常发起，但是返回结果被浏览器拦截了

` 设计思想`：CORS背后的基本思想就是使用自定义的HTTP头部，让 *服务器能声明* 哪些来源可以通过浏览器访问该服务器上的资源，从而决定请求或响应是应该成功还是失败

`解决方案`：CORS，全称是"跨域资源共享"（Cross-origin resource sharing），它允许浏览器向跨源服务器，发出[`XMLHttpRequest`](http://www.ruanyifeng.com/blog/2012/09/xmlhttprequest_level_2.html)请求，从而克服了AJAX只能[同源](http://www.ruanyifeng.com/blog/2016/04/same-origin-policy.html)使用的限制。最实用的方法，是在nginx中配置相应头信息Access-Control-Allow-Origin、Access-Control-Allow-Credentials、Access-Control-Allow-Methods，配置CORS会有两个请求，第一个请求是options预检请求，检查了`Origin`、`Access-Control-Request-Method`和`Access-Control-Request-Headers`字段以后，确认允许跨源请求，就可以做出回应。一旦服务器通过了"预检"请求，以后每次浏览器正常的CORS请求，就都跟简单请求一样，会有一个`Origin`头信息字段。服务器的回应，也都会有一个`Access-Control-Allow-Origin`头信息字段。

```nginx
location / {
    if ($request_method = OPTIONS ) {
        add_header Access-Control-Allow-Origin "http://example.com";
        add_header Access-Control-Allow-Methods "GET, OPTIONS";
        add_header Access-Control-Allow-Headers "Authorization";
        add_header Access-Control-Allow-Credentials "true";
        add_header Content-Length 0;
        add_header Content-Type text/plain;
        return 200;
    }
   	add_header 'Access-Control-Allow-Origin' '*';
    add_header 'Access-Control-Allow-Credentials' 'true';
    add_header 'Access-Control-Allow-Methods' 'POST,OPTIONS';
    proxy_pass http://backend_v2;
}
```

参考：[跨域资源共享 CORS 详解-阮一峰](http://www.ruanyifeng.com/blog/2016/04/cors.html)、[HTTP访问控制(CORS)](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Access_control_CORS) 、[CORS on Nginx](http://enable-cors.org/server_nginx.html) 、[gist nginx cors](https://gist.github.com/pauloricardomg/7084524)

HTTP首部字段类型

`通用首部字段`：请求报文和响应报文两方都会使用的首部

`请求首部字段`：从客户端向服务器端发送请求报文时使用的首部

`响应首部字段` ：从服务器向客户端返回响应报文使用的首部

`实体首部字段`：针对请求报文和响应报文的实体部分使用的首部

HTTP + 加密 + 认证 + 完整性保护 = HTTPS

HTTPS通信

- 步骤一：客户端发送client hello报文开始SSL通信，指定客户端SSL的指定版本、加密组件列表
- 步骤二：服务器端可进行SSL通信、返回Server Hello报文，并返回SSL版本和加密组件列表
- 步骤三：服务器端发送certificate报文，报文中包含公开密钥证书
- 步骤四：服务器端发送server hello done报文通知客户端，最初阶段SSL握手协商结束
- 步骤五：第一阶段结束后，客户端发送Client key Exchange报文作为回应，报文中包含一种Pre-master secret的随机密码串，该报文已用步骤三中的公开密钥进行加密
- 步骤六：客户端接着发送Change cipher spec报文，该报文提示服务器，此报文后的通信会采用Pre-master secret密钥加密
- 步骤七：客户端发送Finished报文，该报文包含连接至今全部报文的整体校验值，这次握手协商是否成功，主要以服务器是否能解密该报文作为判定标准
- 步骤八：服务器同样发送Change Cipher Spec报文
- 步骤九：服务器同样发送Finished报文
- 步骤十：SSL通信完成，开始进行应用层协议通信，即发送HTTP请求

HTTPS使用SSL(secure socket layer)和TLS(transport layer security)两个协议



HTTP keep-alive与websocket有何区别

`区别`：从tcp层来看是一致的，socket都是打开的状态。keep-alive主要是针对浏览器而言，表示浏览器可以复用这个连接从服务器端获取更多的数据，而websocket是双工通信，客户端与服务器端可以在任意时刻进行通信。

socks与http代理协议

`区别`：与HTTP 层代理不同，Socks 代理只是简单地传递数据包，而不必关心是何种应用协议，socks5可参考代码github socks5_go

参考：[socket5 协议学习与实现(一)](http://www.mojidong.com/network/2015/03/07/socket5-1/)、[HTTP 代理原理及实现](https://imququ.com/post/web-proxy.html)

###### websocket协议

websocket是基于http协议的，在握手阶段是一致的，有交集但不是全部

WebSocket的目的就是为了在基础上保证传输的数据量最少

```shell
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: ************==
Sec-WebSocket-Version: **
```

```
响应头
Upgrade：websocket
Connnection: Upgrade
Sec-WebSocket-Accept: ******************
```



###### IO多路复用

`select模型`：

```c
int select (int n, fd_set *readfds, fd_set *writefds,
        fd_set *exceptfds, struct timeval *timeout);
```

select 函数监视的文件描述符分 3 类，分别是 writefds、readfds和 exceptfds。调用后 select 函数会阻塞，直到有描述符就绪（有数据 可读、可写、或者有except），或者超时（timeout 指定等待时间，如果立即返回设为 null 即可）。当 select 函数返回后，通过遍历 fd_set，来找到就绪的描述符。

`poll模型`：

```c
int poll (struct pollfd *fds, unsigned int nfds, int timeout);
//不同与select使用三个位图来表示三个fdset的方式，poll 使用一个 pollfd 的指针实现。
struct pollfd {
    int fd; /* file descriptor */
    short events; /* requested events to watch */
    short revents; /* returned events witnessed */
};
```

pollfd 结构包含了要监视的 event 和发生的 event，不再使用 select “参数-值”传递的方式。同时，pollfd 并没有最大数量限制（但是数量过大后性能也是会下降）。 和 select 函数一样，poll 返回后，需要轮询 pollfd 来获取就绪的描述符。

`epoll模型`：

```c
int epoll_create(int size)；
int epoll_ctl(int epfd, int op, int fd, struct epoll_event *event)；
            typedef union epoll_data {
                void *ptr;
                int fd;
                __uint32_t u32;
                __uint64_t u64;
            } epoll_data_t;

            struct epoll_event {
                __uint32_t events;      /* Epoll events */
                epoll_data_t data;      /* User data variable */
            };

int epoll_wait(int epfd, struct epoll_event * events,
                int maxevents, int timeout);
```

主要是 epoll_create，epoll_ctl 和 epoll_wait 三个函数。epoll_create 函数创建 epoll 文件描述符，参数 size 并不是限制了 epoll 所能监听的描述符最大个数，只是对内核初始分配内部数据结构的一个建议。epoll_ctl 完成对指定描述符 fd 执行 op 操作控制，event 是与 fd 关联的监听事件。op 操作有三种：添加 EPOLL_CTL_ADD，删除 EPOLL_CTL_DEL，修改 EPOLL_CTL_MOD。分别添加、删除和修改对 fd 的监听事件。epoll_wait 等待 epfd 上的 IO 事件，最多返回 maxevents 个事件。

在 select/poll 中，进程只有在调用一定的方法后，内核才对所有监视的文件描述符进行扫描，而 epoll 事先通过 epoll_ctl() 来注册一个文件描述符，一旦基于某个文件描述符就绪时，内核会采用类似 callback 的回调机制，迅速激活这个文件描述符，当进程调用 epoll_wait 时便得到通知。

主要优点有以下几个方面

- 监视的描述符数量不受限制，它所支持的 fd 上限是最大可以打开文件的数目
- IO 的效率不会随着监视 fd 的数量的增长而下降（epoll采用了回调方式监听fd）
- 支持水平触发和边沿触发两种模式（反复通知如果没有通知到，边沿触发只会通知一次）
- mmap 加速内核与用户空间的信息传递（避免频繁的内存拷贝）

###### 惊群

`概念`：多线程/多进程（linux下线程进程也没多大区别）等待同一个socket事件，当这个事件发生时，这些线程/进程被同时唤醒，就是惊群

`nginx中的惊群 `：master进程监听端口号（例如80），所有的nginx worker进程开始用epoll_wait来处理新事件（linux下），也就是各个worker子进程将listenfd加入到自己的epoll中，所以一个新连接来临时，会有多个worker进程在epoll_wait后被唤醒，然后发现自己accept失败。可以通过使用accept_mutex令同一时刻只允许一个nginx worker在自己的epoll中处理监听句柄。



IP地址有网络ID和主机ID组成

[HTTP:80](undefined)
FTP:20/21
SMTP:25
POP3:110
DNS:53
TELNET:23

A类：1.0.0.0--126.0.0.0--用于超大规模网络；
B类：128.0.0.0--191.255.0.0--用于中等规模的网络；
C类：192.0.0.1--223.255.255.0--用于小型的网络；
D类：最高为1110，是多播地址；
E类：最高位是11110，保留在今后使用；

###### OAuth2

![OAuth](/Users/tamchen/Documents/project/material/OAuth.jpg)

第一步：从client 跳转到 服务提供商指定页面，使用如下url，具体参数参看各个服务提供者的要求

https://api.weibo.com/oauth2/authorize?client_id=123050457758183&redirect_uri=http://jianshu.com/callback

第二步：从服务提供商获取access code，通过上一步获取的code

https://api.weibo.com/oauth2/access_token

第三步：访问服务提供商的资源接口，通过access_token

###### 粘包

TCP为了保证可靠传输，尽量减少额外开销（每次发包都要验证），因此采用了流式传输，面向流的传输，相对于面向消息的传输，可以减少发送包的数量。从而减少了额外开销。但是，对于数据传输频繁的程序来讲，使用TCP可能会容易粘包。当然，对接收端的程序来讲，如果机器负荷很重，也会在接收缓冲里粘包。这样，就需要接收端额外拆包，增加了工作量。因此，这个特别适合的是数据要求可靠传输，但是不需要太频繁传输的场合（两次操作间隔100ms，具体是由TCP等待发送间隔决定的，取决于内核中的socket的写法）

而UDP，由于面向的是消息传输，它把所有接收到的消息都挂接到缓冲区的接受队列中，因此，它对于数据的提取分离就更加方便，但是，它没有粘包机制，因此，当发送数据量较小的时候，就会发生数据包有效载荷较小的情况，也会增加多次发送的系统发送开销（系统调用，写硬件等）和接收开销。因此，应该最好设置一个比较合适的数据包的包长，来进行UDP数据的发送。（UDP最大载荷为1472，因此最好能每次传输接近这个数的数据量，这特别适合于视频，音频等大块数据的发送，同时，通过减少握手来保证流媒体的实时性）

出现粘包现象的原因既可能由发送方造成，也可能由接收方造成。

1 发送端需要等缓冲区满才发送出去，造成粘包

2 接收方没能及时地接收缓冲区的包，造成多个包接收



TCP无保护消息边界的解决 针对这个问题，一般有3种解决方案：

 (1)发送固定长度的消息

 (2)把消息的尺寸与消息一块发送

 (3)使用特殊标记来区分消息间隔

 保护消息边界，是指传输协议将数据当做一条独立的消息进行传输，接收端只能接收接收独立消息，每次只能接收一条消息。



建立TCP连接通道、传输数据、断开TCP连接通道

三次握手 粘包（指的是在tcp流式传输出现的无法区分消息边界的情况，UDP不会粘包是因为它有消息边界）

超时重传、快速重传、流量控制、拥塞控制（慢启动、拥塞避免、拥塞发生、快速恢复）

四次挥手 可以看做两次断开连接

SYN包、ACK包、FIN包

###### 负载均衡

所谓四层就是基于IP+端口的负载均衡（ TCP层）；七层就是基于URL等应用层信息的负载均衡；

四层负载均衡：主要通过报文中的目标地址和端口，再加上负载均衡设备设置的服务器选择方式，决定最终选择的内部服务器。

七层负载均衡：主要通过报文中的真正有意义的应用层内容，再加上负载均衡设备设置的服务器选择方式，决定最终选择的内部服务器

参考书籍

[《图解TCP/IP协议》](http://www.jianshu.com/p/de262cfbb4ef)

[《UNIX网络编程 卷1：套接字联网API - 读书笔记》](http://cstdlib.com/tech/2014/10/09/read-unix-network-programming-1/)

[谈一谈网络编程学习经验](https://cloud.github.com/downloads/chenshuo/documents/LearningNetworkProgramming.pdf)

[APUE Advanced Programming Unix Environment](http://dirtysalt.github.io/apue.html#orgheadline29)

[《TCP的那些事》陈皓](http://coolshell.cn/articles/11564.html)

参考资料

- [openresty best practices](https://www.gitbook.com/book/moonbingbing/openresty-best-practices)
- [nginx open source architecture Andrew Alexeev 中文译文](http://www.ituring.com.cn/article/4436)
- [理解OAuth 2.0](http://www.ruanyifeng.com/blog/2014/05/oauth_2_0.html)
- [使用 OAuth 2.0 访问豆瓣 API](https://developers.douban.com/wiki/?title=oauth2)
- [tcp粘包分析](http://lib.csdn.net/article/computernetworks/17066)