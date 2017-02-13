#### nginx架构

###### nginx进程模型

`进程模型`：ngnix在Unix系统以daemon方式运行，后台进程包括一个master进程和多个worker进程

`进程任务`：master进程任务，接收来自外界的信号，向各worker进程发送信号，监控worker进程的运行状态。当worker进程退出后(异常情况下)，会自动重新启动新的worker进程，worker进程任务，accept连接请求，处理request

`进程管理`：可以通过直接给master进程发送信号来管理worker，例如kill -HUP pid，则是告诉master，重新加载配置文件，启动新的worker进程，让老的worker退休，从容地重启nginx，或者通过命令行方式，./nginx -s reload，该命令会新建一个nginx进程，接收解析命令，后向master通信

###### nginx事件模型

`建立listenfd`：在master进程里面，先建立好需要listen的socket（listenfd）之后，然后再fork出多个worker进程。所有worker进程的listenfd会在新连接到来时变得可读，

`抢占accept_mutex`：为保证只有一个进程处理该连接，所有worker进程在注册listenfd读事件前抢accept_mutex，抢到互斥锁的那个进程注册listenfd读事件，在读事件里调用accept接受该连接，此后就由该worker负责后续的读请求、解析请求、处理请求、产生数据、返回客户端、断开连接。

`异步非阻塞处理`：使用select/poll/epoll/kqueue这样的系统调用。它们提供了一种机制，让你可以同时监控多个事件，调用他们是阻塞的，但可以设置超时时间，在超时时间之内，如果有事件准备好了，就返回（处理多连接，避免上下文切换）

`epoll事件模型`：在epoll中，当事件没准备好时，放到epoll里面，事件准备好了，我们就去读写，当读写返回EAGAIN时，我们将它再次加入到epoll里面，让出线程处理。只要有事件准备好了，就去处理它，只有当所有事件都没准备好时，才在epoll里面等着。这样，就可以实现在一个线程中循环处理多个准备好的事件了。

`与多线程模型对比`：不需要创建线程，每个请求内存很小，并发数多也不会导致开销很大的上下文切换

###### nginx connection概念

`基本结构`：connection包含连接的socket、读事件、写事件，可以使用connection处理建立连接，发送数据，接收数据

`nginx处理tcp连接`：master进程启动，创建socket，设置addrreuse等选项，绑定到指定的ip地址端口，再listen，fork多个worker子进程，子进程会竞争accept新的连接，此后客户端就与服务器进行三次握手连接，worker子进程会将该socket封装成ngx_connection_t结构体

`worker_connections`：nginx通过设置worker_connectons来设置每个进程支持的最大连接数，nginx在实现上设置了一个worker_connections大小ngx_connection_t结构的数组，同时通过free_connections链表保存所有空闲的ngx_connection_t，每次新建一个连接就从free_connections拿一个，用完放回，worker_connections指的是每个worker的最大连接数，而不是整个nginx的最大连接数

`nginx竞争accept处理`：nginx使用一个叫ngx_accept_disabled的变量来控制是否去竞争accept_mutex锁。这个值是nginx单进程的所有连接总数的八分之一，减去剩下的空闲连接数量，当ngx_accept_disabled大于0时，不会去尝试获取accept_mutex锁，并且将ngx_accept_disabled减1，ngx_accept_disabled越大让出连接就越频繁，其他进程获取accept_mutex机会越大

###### nginx request概念

`基本结构`：ngx_http_request_t是对一个http请求的封装。 一个http请求，包含请求行、请求头、请求体、响应行、响应头、响应体。

`request处理过程`：处理过程为请求头、请求体，响应头，响应体，ngx_http_init_request —>设置读事件为ngx_http_process_request_line(处理请求行）nginx采用状态机对请求行分析，并将method四个字符转化为整形一次比较减少cpu指令数，—>ngx_http_process_request_headers —>ngx_http_process_request，具体可以查看nginx开发从入门到精通

###### 长连接keep-alive

`定义`：在一个连接上面执行多个请求的，这就是所谓的长连接

`请求头`：请求头必须指定content-length表明请求体大小，否则服务返回400错误

`响应头`：根据http协议版本不同content-length规定也不同，如果指定了大小，则客户端按照指定大小接收数据，如果没有指定，则在http1.0，客户端直接关闭连接，而在http1.1中客户端则一直等待服务器响应

`客户端请求头connection`：如果客户端请求头connection为close，则表示客户端关闭keep-alive，如果为keep-alive则需要打开长连接，如果没有指定，则按照协议版本，nginx在响应完后会设置keep-alive，同时会设置keep-alive timeout，对于请求量大的客户端来说，如果设置了keep-alive close，则会产生大量的time-wait

###### nginx模块化

`模块化`：nginx的内部结构是由核心部分和一系列的功能模块所组成，每个模块实现特定的功能，当一个请求到达时需要经过nginx模块链中的部分或全部模块，主要分为几类event module、phase module、output module、upstream、load-balancer

###### nginx处理请求

`worker处理流程`：

1. 操作系统提供的机制（例如epoll, kqueue等）产生相关的事件。
2. 接收和处理这些事件，如是接受到数据，则产生更高层的request对象。
3. 处理request的header和body。
4. 产生响应，并发送回客户端。
5. 完成request的处理。
6. 重新初始化定时器及其他事件。

###### phase handler

`阶段处理模块`：该模块包含若干request处理阶段的handler，主要执行以下任务

1. 获取location配置。
2. 产生适当的响应。
3. 发送response header。
4. 发送response body。

###### nginx配置

`server select`：by port —> by address —>by server_name(host)—>find virtual server names?

—> found core server conf || not found find in virtual server regex names

—> core server conf —> end

`location match`：location指令格式如下 location optional_modifier location_match，例如 location [=|~|~*|^~] /uri/ { … }

optional_modifier 有以下几种配置

- (none)，如果没有optional_modifier，location被解析为前缀匹配。这意味着给定的location需要与请求的URI开始部分是匹配的。
- =，这表示给定的location需要请求的URI完全匹配，也就是精确匹配。
- ~，这表示大小写敏感的正则表达式匹配。
- ~*，这表示大小写不敏感的正则表达式匹配。
- ^*，这表示最佳的非正则表达式匹配。
- /，通用匹配，所有请求都会匹配

location匹配过程

- 首先匹配 =，如果匹配则结束
- 然后匹配 ^*，如果匹配则结束
- 其次按照文件顺序进行前缀匹配，匹配则保存最长匹配项
- 按照正则表达式匹配，如果匹配则匹配结束
- 正则匹配不成功，采用前缀匹配保存的最长匹配项
- 都没有匹配到，则返回not found

`URL Rewrite`：使用nginx提供的全局变量或自己设置的变量，结合正则表达式和标志位实现url重写以及重定向，rewrite只在server{}、location{}、if{}块中起作用，语法rewrite regex replacement [flag];

flag标志位

- `last` : 相当于Apache的[L]标记，表示完成rewrite
- `break` : 停止执行当前虚拟主机的后续rewrite指令集
- `redirect` : 返回302临时重定向，地址栏会显示跳转后的地址
- `permanent` : 返回301永久重定向，地址栏会显示跳转后的地址

###### 反向代理与正向代理

`反向代理`：反向代理（Reverse Proxy）方式是指用代理服务器来接受 internet 上的连接请求，然后将请求转发给内部网络上的服务器，并将从服务器上得到的结果返回给 internet 上请求连接的客户端，此时代理服务器对外就表现为一个反向代理服务器，在ngnix中主要使用proxy_pass将请求反向代理到其他服务器，在默认情况下不会转发原始请求中host头部，如果需要转发，可使用proxy_set_header Host $host，如果需要转发默认host，可配置$proxy_host，具体用法可以到[Nginx 官网](http://nginx.org/en/docs/http/ngx_http_proxy_module.html)查看

`正向代理`：正向代理就像一个跳板，常见的翻墙工具，游戏代理都是通过正向代理实现

###### upstream 负载均衡的用法

`基本用法`：

```nginx
# 在http节点下，配置upstream
upstream favtomcat {
       server 10.0.6.108:7080;
       server 10.0.0.85:8980;
}
# 在某个server节点的location节点下配置
location / {
        root   html;
        index  index.html index.htm;
        proxy_pass http://favtomcat;
}
```

`分配策略`：weight权重，ip_hash（按IPhash结果），fair（第三方）url_hash(第三方)

`为每个设备设置状态值`：down、weight、max_fails(允许请求失败的次数)、fails_timeout(max_fails次失败后，暂停的时间)、backup(表明这是台备份机子)

###### 可用的全局变量

- $args
- $content_length
- $content_type
- $document_root
- $document_uri
- $host
- $http_user_agent
- $http_cookie
- $limit_rate
- $request_body_file
- $request_method
- $remote_addr
- $remote_port
- $remote_user
- $request_filename
- $request_uri
- $query_string
- $scheme
- $server_protocol
- $server_addr
- $server_name
- $server_port
- $uri

如何用golang实现这样的一个nginx系统

###### nginx 使用多进程而不是多线程

- 进程间不共享资源，不需要加锁，省去锁带来的开销
- 采用独立进程，保证服务的可用性，一个work进程挂了不影响其他进程
- 多线程带来的上下文切换，造成巨大的CPU开销

###### 参考资料

- [http://aosabook.org/en/nginx.html](http://aosabook.org/en/nginx.html)
- [nginx开发从入门到精通](http://tengine.taobao.org/book/)
- [Nginxserver和location选择策略](http://www.php101.cn/2015/11/12/Nginxserver%E5%92%8Clocation%E9%80%89%E6%8B%A9%E7%AD%96%E7%95%A5/)
- [nginx-internals](http://www.slideshare.net/joshzhu/nginx-internals)
- [nginx配置location总结及rewrite规则写法](http://seanlook.com/2015/05/17/nginx-location-rewrite/)
- [nginx best practices](https://moonbingbing.gitbooks.io/openresty-best-practices/content/ngx/pitfalls_and_common_mistakes.html)
- [书籍 nginx essential]()
- ​


