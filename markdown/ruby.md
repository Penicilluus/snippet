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