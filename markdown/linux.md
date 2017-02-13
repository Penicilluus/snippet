###### socket编程

​	A *socket* is an abstraction of a communication endpoint and it allows your program to communicate with another program using file descriptors.

​	The *socket pair* for a TCP connection is a 4-tuple that identifies two endpoints of the TCP connection: the local IP address, local port, foreign IP address, and foreign port. A socket pair uniquely identifies every TCP connection on a network. The two values that identify each endpoint, an IP address and a port number, are often called a *socket*.[1](https://ruslanspivak.com/lsbaws-part3/#fn:1)

​	The TCP/IP stack within the kernel automatically assigns the local IP address and the local port when the client calls *connect*. The local port is called an *ephemeral port*, 

socket 创建流程

######  进程

​	In UNIX, every user process also has a parent that, in turn, has its own process ID called parent process ID, or PPID for short.

###### 文件描述符

​	So what is a file descriptor? A *file descriptor* is a non-negative integer that the kernel returns to a process when it opens an existing file, creates a new file or when it creates a new socket. 

​	By default, UNIX shells assign file descriptor 0 to the standard input of a process, file descriptor 1 to the standard output of the process and file descriptor 2 to the standard error.

###### 子进程创建 fork

​	The most important point to understand about [fork()](https://docs.python.org/2.7/library/os.html#os.fork) is that you call *fork* once but it returns twice: once in the parent process and once in the child process. When you fork a new process the process ID returned to the child process is 0. When the *fork* returns in the parent process it returns the child’s PID.

When a parent forks a new child, the child process gets a copy of the parent’s file descriptors:

​	The kernel uses descriptor reference counts to decide whether to close a socket or not.  It closes the socket only when its descriptor reference count becomes 0. When your server creates a child process, the child gets the copy of the parent’s file descriptors and the kernel increments the reference counts for those descriptors. 

###### 僵尸进程 zombie

​	A *zombie* is a process that has terminated, but its parent has not *waited* for it and has not received its termination status yet. When a child process exits before its parent, the kernel turns the child process into a zombie and stores some information about the process for its parent process to retrieve later. The information stored is usually the process ID, the process termination status, and the resource usage by the process.

###### Unix I/O模型

###### linux命令

查看进程打开的文件lsof -c programe-name 或lsof -p $PID

查看监听端口netstat -anp 所有监听端口及对应的进程

netstat -an | grep LISTEN  显示所有socket     -n 以**网络**IP地址代替名称，显示出**网络**连接情形



select/poll/epoll



- 写一条crontab配置, 每周六，日重启Nginx, 并将执行结果写入`/dev/null`



###### ubuntu实用命令

配置Ubuntu开机自启动脚本

```shell
# Ubuntu开机自启动脚本位置 /etc/rc.local
sudo su - kaiwan -c "sudo nginx"
exit 0
# 
# update-rc.d docker defaults
Adding system startup for /etc/init.d/docker ...
/etc/rc0.d/K20docker -> ../init.d/docker
/etc/rc1.d/K20docker -> ../init.d/docker
/etc/rc6.d/K20docker -> ../init.d/docker
/etc/rc2.d/S20docker -> ../init.d/docker
/etc/rc3.d/S20docker -> ../init.d/docker
/etc/rc4.d/S20docker -> ../init.d/docker
/etc/rc5.d/S20docker -> ../init.d/docker
```

crontab任务

```shell
crontab [-u user] file crontab [-u user] [ -e | -l | -r ]

crontab -e #进入编辑定时命令页面
# *  *  * * * sh /home/admin/test/test.sh
# 每分钟执行test.sh脚本

#基本格式 :
#*　　*　　*　　*　　*　　command
#分　时　日　月　周　命令
#/etc/crontab 
1、每分钟执行一次            
*  *  *  *  * 

2、每隔一小时执行一次        
00  *  *  *  * 
or
* */1 * * *  (/表示频率)

3、每小时的15和30分各执行一次 
15,45 * * * * （,表示并列）

4、在每天上午 8- 11时中间每小时 15 ，45分各执行一次
15,45 8-11 * * * command （-表示范围）

5、每个星期一的上午8点到11点的第3和第15分钟执行
3,15 8-11 * * 1 command

6、每隔两天的上午8点到11点的第3和第15分钟执行
3,15 8-11 */2 * * command
```



[Linux 下执行定时任务 crontab 命令详解](https://segmentfault.com/a/1190000002628040)



什么是大小端(0x2211)

- **大端字节序**：高位字节在前，低位字节在后，这是人类读写数值的方法。
- **小端字节序**：低位字节在前，高位字节在后，即以`0x1122`形式储存。

为什么会有小端字节序？

答案是，计算机电路先处理低位字节，效率比较高，因为计算都是从低位开始的。所以，计算机的内部处理都是小端字节序

如果是大端字节序，先读到的就是高位字节，后读到的就是低位字节。小端字节序正好相反。

"只有读取的时候，才必须区分字节序，其他情况都不用考虑。"



rednax:虽然x86是小端的，不过也有很多CPU是大端的（比如powerpc）。另外小端虽然对加法器比较友好，但除法还是大端更合适，所以这种取舍其实只是一些历史问题吧

小端模式 ：强制转换数据不需要调整字节内容，1、2、4字节的存储方式一样。
大端模式 ：符号位的判定固定为第一个字节，容易判断正负。

小端模式更适合系统内部，大端模式更适合网络数据传递，加上一些历史引领的原因，
导致现在两种字节序方式并存。

某些机器选择在存储器中按照从最低有效字节到最高有效字节的顺序存储对象，而另一些机器则按照从最高有效字节到最低有效字节的顺序存储。前一种规则——最低有效字节在最前面的方式，称为**小端法(little endian)**。后一种规则——最高有效字节在最前面的方式，称为**大端法(big endian)**。许多比较新的微处理器使用**双端法(bi-edian)**，也就是说可以把它们配置成作为大端或者小端的机器运行。



###### ubuntu执行脚本的几种方法

`source命令`：在当前bash环境下，执行文件中的命令，该文件可以无执行权限，通常用 . 替代，在文件中的变量设置会对当前的bash环境产生影响

`sh bash命令`：打开一个子shell（子进程）读取文件中的命令并执行，该文件可以无执行权限，在文件中的变量设置对当前bash环境无影响

`./命令`：打开一个子shell 执行文件中的内容，需要执行权限，由于打开了另外一个命令解释器，所以对当前bash环境无影响

